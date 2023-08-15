package handlers

import (
	"TESTGO/pkg/api/errors"
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

func TrueIDSubscripber(client external.TrueIDSubscripberAPI, c *gin.Context, db *gorm.DB) {
	if ssoidValue, exists := c.Get("ssoid"); exists {
		ssoid, ok := ssoidValue.(string)
		if !ok {
			fmt.Println("ssoid is not string")
			c.JSON(errors.ErrValidationInputSSOID.Status, errors.ErrValidationInputSSOID)
		}
		subscription, resp, err := client.GetSubscribers(ssoid)
		if err != nil {
			var resultMap map[string]interface{}
			err := json.Unmarshal([]byte(resp.String()), &resultMap)
			if err != nil {
				c.JSON(errors.ErrParseJSON.Status, errors.ErrParseJSON)
				return
			}
			c.JSON(resp.StatusCode(), resultMap)
			return
		}
		c.JSON(resp.StatusCode(), subscription)
	} else {
		c.JSON(errors.ErrValidationInputSSOID.Status, errors.ErrValidationInputSSOID)
	}
}
