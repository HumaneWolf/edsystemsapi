package api

import (
	"strings"

	"github.com/gin-gonic/gin"
	"humanewolf.com/ed/systemapi/config"
)

// CheckAccessControl checks the config to determine if access control should be active, and verifies token if so.
func CheckAccessControl() gin.HandlerFunc {
	return func(c *gin.Context) {
		cfg := config.LoadConfig()

		if cfg.AccessControl.RequireAccessToken {
			token, hasToken := getAccessToken(c)

			if !hasToken || token != cfg.AccessControl.AccessToken {
				c.AbortWithStatusJSON(403, gin.H{"error": "Missing or incorrect bearer token."})
				return
			}
		}

		c.Next()
	}
}

func getAccessToken(c *gin.Context) (string, bool) {
	authHeaders := c.Request.Header["Authorization"]
	if len(authHeaders) < 1 {
		return "", false
	}

	headerParts := strings.Split(authHeaders[0], " ")
	if len(headerParts) < 2 {
		return "", false
	}

	tokenType := headerParts[0]
	token := headerParts[1]

	if strings.ToLower(tokenType) != "bearer" {
		return "", false
	}

	return token, true
}
