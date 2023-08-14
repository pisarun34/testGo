package api

import (
	"TESTGO/api/handlers"
	"TESTGO/middlewares"
	"TESTGO/pkg/database"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(r *gin.Engine, redis database.RedisClientInterface, db *gorm.DB) {
	baseGroup := r.Group("/")
	baseGroup.Use(middlewares.AuthTrueID())

	setupSeeksterRoutes(baseGroup, redis, db)
	setupTrueIDRoutes(baseGroup, redis, db)
}

func setupSeeksterRoutes(baseGroup *gin.RouterGroup, redis database.RedisClientInterface, db *gorm.DB) {
	seeksterClientInstance := handlers.NewSeeksterClient()
	seeksterGroup := baseGroup.Group("/seekster")
	seeksterGroup.Use(middlewares.AuthSeekster())
	{
		seeksterGroup.POST("/signin", func(c *gin.Context) {
			handlers.SeeksterSignin(seeksterClientInstance, c, redis, db)
		})
		seeksterGroup.POST("/signup", func(c *gin.Context) {
			handlers.SeeksterSignup(seeksterClientInstance, c, db)
		})
		seeksterGroup.POST("/insertuser", func(c *gin.Context) {
			handlers.InsertSeeksterUser(seeksterClientInstance, c, db)
		})
	}
}

func setupTrueIDRoutes(baseGroup *gin.RouterGroup, redis database.RedisClientInterface, db *gorm.DB) {
	trueidClientInstance := handlers.NewTrueIDClient()
	trueIDGroup := baseGroup.Group("/trueid")
	trueIDGroup.Use(middlewares.AuthTrueID())
	{
		trueIDGroup.POST("/subscriber", func(c *gin.Context) {
			handlers.TrueIDSubscripber(trueidClientInstance, c, db)
		})

	}
}
