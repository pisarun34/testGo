package middlewares

import (
	"TESTGO/pkg/api/errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				apiErr, ok := err.Err.(*errors.APIError)
				if ok {
					c.JSON(apiErr.Status, apiErr)
					return
				}

				externalAPIErr, ok := err.Err.(*errors.ExternalAPIError)
				if ok {
					c.JSON(externalAPIErr.Status, externalAPIErr.ErrorResponse)
					return
				}
			}

			// ถ้าเป็นข้อผิดพลาดอื่นที่ไม่ได้กำหนดรูปแบบ, ส่ง HTTP 500
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "Internal Server Error",
			})
		}
		c.Next()
	}
}
