package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Logs() gin.HandlerFunc {
	return func(c *gin.Context) {

		// ประมวลผล request
		c.Next()

		// หลังจากประมวลผล request เสร็จ
		status := c.Writer.Status()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		if raw != "" {
			path = path + "?" + raw
		}

		// ข้อมูลที่จะ log
		logrus.WithFields(logrus.Fields{
			"status":     status,
			"method":     c.Request.Method,
			"path":       path,
			"ip":         c.ClientIP(),
			"user-agent": c.Request.UserAgent(),
			"errors":     c.Errors.String(),
		}).Info("request completed")

		c.Next()
	}
}
