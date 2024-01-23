package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// cors is a middleware that adds CORS headers to requests.
//
// It allows requests from the specified origin.
//
// It also handles preflight requests via the OPTIONS method. In this case,
// the request is aborted with a 200 OK status.
func cors(allowOrigin string) gin.HandlerFunc {
	const (
		allowCredentials = "true"
		allowMethods     = "OPTIONS, POST, GET, PUT, DELETE"
		allowHeaders     = "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, " +
			"Authorization, Accept, Origin, Cache-Control, X-Requested-With"
	)

	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", allowOrigin)
		c.Header("Access-Control-Allow-Credentials", allowCredentials)
		c.Header("Access-Control-Allow-Methods", allowMethods)
		c.Header("Access-Control-Allow-Headers", allowHeaders)

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusOK)
		} else {
			c.Next()
		}
	}
}
