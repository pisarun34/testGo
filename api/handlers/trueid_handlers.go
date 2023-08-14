package handlers

import (
	"TESTGO/models"
	"TESTGO/pkg/external"
	"TESTGO/pkg/external/trueidpartner"
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewTrueIDClient() external.TrueIDSubscripberAPI {
	return trueidpartner.NewClient()
}

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func TrueIDSubscripber(client external.TrueIDSubscripberAPI, c *gin.Context, db *gorm.DB) {
	if ssoidValue, exists := c.Get("ssoid"); exists {
		ssoid, ok := ssoidValue.(string)
		if !ok {
			fmt.Println("ssoid is not string")
			c.JSON(models.ErrValidationInputSSOID.Status, models.ErrValidationInputSSOID)
		}
		subscription, resp, err := client.GetSubscribers(ssoid)
		if err != nil {
			var resultMap map[string]interface{}
			err := json.Unmarshal([]byte(resp.String()), &resultMap)
			if err != nil {
				c.JSON(models.ErrParseJSON.Status, models.ErrParseJSON)
				return
			}
			c.JSON(resp.StatusCode(), resultMap)
			return
		}
		c.JSON(resp.StatusCode(), subscription)
	} else {
		c.JSON(models.ErrValidationInputSSOID.Status, models.ErrValidationInputSSOID)
	}
}
