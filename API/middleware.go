package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Clayagiffeb/Simple_Bank/token"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey     = "authorization"
	authorizationHeaderBrearer = "bearer" // for simplicity
	authorizationPayloadKey    = "authorization_payload"
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is empty")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		fields := strings.Fields(authorizationHeader) // split authorization header
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationHeaderBrearer {
			err := fmt.Errorf("invalid type authorization %s", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		// now we will check the access token

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		// valid token
		ctx.Set(authorizationPayloadKey, payload) // a key-value pair
		ctx.Next()                                // forward it to the next handler
	}
}
