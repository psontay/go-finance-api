package api

import (
	"SimpleBank/token"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

// RateLimiterMiddleware creates a middleware that limits requests by User IP or Username
func RateLimiterMiddleware(redisClient *redis.Client, limit int, window time.Duration) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var identifier string
		authPayload, exists := ctx.Get(authorizationPayloadKey)
		if exists {
			payload, ok := authPayload.(*token.Payload)
			if ok {
				identifier = payload.Username
			}
		}

		if identifier == "" {
			identifier = ctx.ClientIP()
		}

		key := fmt.Sprintf("rate_limit:%s:%s", ctx.FullPath(), identifier)

		pipe := redisClient.TxPipeline()
		incr := pipe.Incr(ctx, key)
		pipe.Expire(ctx, key, window)
		_, err := pipe.Exec(ctx)

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to check rate limit: %v", err)))
			return
		}

		if incr.Val() > int64(limit) {
			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded, please try again later"})
			return
		}

		ctx.Next()
	}
}

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// is user has header authorization?
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		// check format of header authorization
		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		// check type of first field, must equal to "bearer"
		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		// get access token
		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}

func roleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authPayload, exists := ctx.Get(authorizationPayloadKey)
		if !exists {
			err := errors.New("authentication payload not found")
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		payload, ok := authPayload.(*token.Payload)
		if !ok {
			err := errors.New("invalid authentication payload type")
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		hasAccountedRole := false
		for _, role := range allowedRoles {
			if payload.Role == role {
				hasAccountedRole = true
				break
			}
		}
		if !hasAccountedRole {
			err := errors.New("user does not have permission to access this resource")
			ctx.AbortWithStatusJSON(http.StatusForbidden, errorResponse(err))
			return
		}
		ctx.Next()
	}
}
