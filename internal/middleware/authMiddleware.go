package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type authZContextKey int

const (
	authorizationHeaderKey                  = "authorization"
	authorizationTypeBearer                 = "bearer"
	AuthZContextKey         authZContextKey = iota
)

type SessionChecker interface {
	CheckSession(sid uuid.UUID) bool
}

func tokenFromHeader(c *gin.Context) (string, error) {
	authorizationHeader := c.GetHeader(authorizationHeaderKey)
	if len(authorizationHeader) == 0 {
		return "", errors.New("authorization header is not provided")
	}

	fields := strings.Fields(authorizationHeader)
	if len(fields) < 2 {
		return "", errors.New("invalid authorization header format")
	}

	authorizationType := strings.ToLower(fields[0])
	if authorizationType != authorizationTypeBearer {
		return "", fmt.Errorf("unsupported authorization type %s", authorizationType)
	}
	accessToken := fields[1]
	return accessToken, nil
}

func AuthMiddleware(authServ SessionChecker) gin.HandlerFunc {
	return func(c *gin.Context) {

		accessToken, err := tokenFromHeader(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		sid, err := uuid.Parse(accessToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		good := authServ.CheckSession(sid)
		fmt.Println("CheckSession", good)
		if !good {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{})
			return
		}

		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, AuthZContextKey, sid)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
