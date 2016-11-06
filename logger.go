/*

src: http://codegists.com/snippet/go/gin-loggergo_logrusorgru_go

// Creates a router without any middleware by default
    r := gin.New()

    r.Use(gin.Recovery())
    r.Use(Logger()) // use this logger without colors

*/
package main

import (
	"github.com/gin-gonic/gin"

	"fmt"
	"os"
	"time"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		// Process request
		c.Next()
		// Stop timer
		end := time.Now()
		latency := end.Sub(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		comment := c.Errors.ByType(gin.ErrorTypePrivate).String()
		fmt.Fprintf(os.Stdout,
			"[GIN] %v | %3d | %13v | %s | %-7s %s\n%s",
			end.Format("2006/01/02 - 15:04:05"),
			statusCode,
			latency,
			clientIP,
			method,
			path,
			comment,
		)
	}
}
