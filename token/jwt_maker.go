package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const minSecretKeySize = 32

// make jwt
type JWTMaker struct {
	secretKey string
}

// put Payload in here
type jwtCustomClaims struct {
	Payload *Payload `json:"payload"`
	jwt.RegisteredClaims
}

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d char", minSecretKeySize)
	}
	return &JWTMaker{secretKey}, nil
}

func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", payload, err
	}

	// put payload in claims
	claims := jwtCustomClaims{
		Payload: payload,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "elvis",
			Subject:   payload.Username,
			ExpiresAt: jwt.NewNumericDate(payload.ExpiredAt),
			IssuedAt:  jwt.NewNumericDate(payload.IssuedAt),
			ID:        payload.ID.String(), // convert uuid to string
		},
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := jwtToken.SignedString([]byte(maker.secretKey))
	return signedToken, payload, nil
}

func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	// extract token
	jwtToken, err := jwt.ParseWithClaims(token, &jwtCustomClaims{}, func(token *jwt.Token) (any, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(maker.secretKey), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
	}
	claims, ok := jwtToken.Claims.(*jwtCustomClaims)
	if !ok {
		return nil, ErrInvalidToken
	}
	err = claims.Payload.Valid()
	if err != nil {
		return nil, err
	}
	return claims.Payload, nil
}
