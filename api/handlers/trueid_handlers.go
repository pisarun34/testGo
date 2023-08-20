package handlers

import (
	"TESTGO/middlewares"
	"TESTGO/pkg/external"
	"TESTGO/pkg/external/trueidpartner"
	"TESTGO/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewTrueIDClient() external.TrueIDSubscripberAPI {
	return trueidpartner.NewClient()
}

func SubscripberAPICall(client external.TrueIDSubscripberAPI, ssoid string) (*trueidpartner.SubscribersResponse, error) {
	subscription, resp, err := client.GetSubscribers(ssoid)
	if errHandle := utils.HandleAPIError(resp, err); errHandle != nil {
		return nil, errHandle
	}
	return subscription, nil
}

func TrueIDSubscripber(client external.TrueIDSubscripberAPI, c *gin.Context, db *gorm.DB) {
	ssoid, err := middlewares.CheckAndExtractSSOID(c)
	if err != nil {
		c.Error(err)
		return
	}
	subscribers, err := SubscripberAPICall(client, ssoid)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, subscribers)

}
