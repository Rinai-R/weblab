package middleware

import (
	"net/http"
	"strings"
	"weblab/internal/utils"

	"github.com/gin-gonic/gin"
)

const ContextUserIDKey = "userID"

func JWTAuth(jwtMgr *utils.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.Fail(c, http.StatusUnauthorized, "missing authorization header")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			utils.Fail(c, http.StatusUnauthorized, "invalid authorization header")
			c.Abort()
			return
		}

		userID, err := jwtMgr.Parse(parts[1])
		if err != nil {
			utils.Fail(c, http.StatusUnauthorized, "invalid token")
			c.Abort()
			return
		}

		c.Set(ContextUserIDKey, userID)
		c.Next()
	}
}
