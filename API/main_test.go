package api

import (
	"os"
	"testing"
	"time"

	db "github.com/Clayagiffeb/Simple_Bank/db/sqlc"
	"github.com/Clayagiffeb/Simple_Bank/util"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	config := util.Config{
		AccessTokenDuration: time.Minute,
	}
	server, err := NewServer(config, store)
	require.NoError(t, err)

	return server

}
func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
