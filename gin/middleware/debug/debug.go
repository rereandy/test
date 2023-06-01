package debug

import (
	"fmt"
	"time"

	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
)

// DebugMiddleWare 执行其它中间件完成后再执行这个DebugMiddleWare
func DebugMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		req := c.Request
		startTime := time.Now()
		fmt.Printf("[Debug Middleware] - Request: %v\n", req)
		c.Next()
		fmt.Printf("[Debug Middleware] - Response: %v\n", c.Writer.Status())
		fmt.Printf("Request API: %v, cost time: %v\n", c.Request.URL, time.Since(startTime))
	}
}

func TimeoutMiddleWare() gin.HandlerFunc {
	return timeout.New(timeout.WithTimeout(3*time.Second),
		timeout.WithHandler(func(c *gin.Context) {
			time.Sleep(5 * time.Second)
			c.Next()
		}),
		timeout.WithResponse(func(c *gin.Context) {
			c.JSON(200, gin.H{
				"msg":  "timeout",
				"code": 0,
			})
		}))
}

// ReportMiddleWare 先执行完这个中间件
func ReportMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		req := c.Request
		startTime := time.Now()
		fmt.Printf("[Debug Middleware] - Request: %v\n", req)
		fmt.Printf("[Debug Middleware] - Response: %v\n", c.Writer.Status())
		fmt.Printf("Request API: %v, cost time: %v\n", c.Request.URL, time.Since(startTime))

	}
}
