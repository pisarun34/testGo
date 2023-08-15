package handlers

import (
	"TESTGO/pkg/api/errors"
	"TESTGO/pkg/database"
	"TESTGO/pkg/database/dtos"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetPackagesList(c *gin.Context, redis database.RedisClientInterface, db *gorm.DB) {
	var packages []dtos.Packages
	if err := db.Preload("SubscriptionPackages").
		Preload("SubscriptionPackages.Price").
		Find(&packages).Error; err != nil {
		c.JSON(errors.ErrDatabase.Status, errors.ErrDatabase)
		return
	}
	c.JSON(200, packages)
}
