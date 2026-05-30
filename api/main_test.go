package api

import (
	database "SimpleBank/database/sqlc"
	"SimpleBank/util"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store database.Store) *Server {
	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	// For unit tests, we don't need a real Redis connection if we're not testing it.
	// But Server needs a valid *redis.Client struct. We can just use an empty one,
	// or ideally use a mock like miniredis. For now, an un-connected client is fine 
	// as long as our tests don't actually trigger Redis commands.
	var redisClient *redis.Client

	server, err := NewServer(config, store, redisClient)
	require.NoError(t, err)

	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
