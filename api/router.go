package api

import (
	"TESTGO/api/handlers"
	"TESTGO/middlewares"
	"TESTGO/pkg/database"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRouter is a function that setup router
func SetupRouter(r *gin.Engine, redis database.RedisClientInterface, db *gorm.DB) {
	baseGroup := r.Group("/")
	// use middleware to extract ssoid from jwt to context
	baseGroup.Use(middlewares.AuthTrueID(), middlewares.ErrorHandler(), middlewares.Logs())
	// setup seekster and TrueID routes
	setupSeeksterRoutes(baseGroup, redis, db)
	setupTrueIDRoutes(baseGroup, redis, db)
	setupTruexRoutes(baseGroup, redis, db)
}

// setupSeeksterRoutes is a function that setup seekster routes
func setupSeeksterRoutes(baseGroup *gin.RouterGroup, redis database.RedisClientInterface, db *gorm.DB) {
	seeksterClientInstance := handlers.NewSeeksterClient()
	seeksterGroup := baseGroup.Group("/seekster")

	seeksterGroup.POST("/signin", func(c *gin.Context) {
		handlers.SeeksterSignin(seeksterClientInstance, c, redis, db)
	})
	seeksterGroup.POST("/signup", func(c *gin.Context) {
		handlers.SeeksterSignup(seeksterClientInstance, c, redis, db)
	})
	seeksterGroup.POST("/insertuser", func(c *gin.Context) {
		handlers.InsertSeeksterUser(seeksterClientInstance, c, db)
	})
	//middlewares.AuthSeekster(seeksterClientInstance, redis, db)
	seeksterGroup.Use(handlers.AuthSeekster(seeksterClientInstance, redis, db))
	{
		// Protected routes
		seeksterGroup.GET("/services", func(c *gin.Context) {
			handlers.GetServiceList(seeksterClientInstance, c, redis, db)
		})
	}
}

// setupTrueIDRoutes is a function that setup TrueID routes
func setupTrueIDRoutes(baseGroup *gin.RouterGroup, redis database.RedisClientInterface, db *gorm.DB) {
	trueidClientInstance := handlers.NewTrueIDClient()
	trueIDGroup := baseGroup.Group("/trueid")
	trueIDGroup.Use()
	{
		trueIDGroup.POST("/subscriber", func(c *gin.Context) {
			handlers.TrueIDSubscripber(trueidClientInstance, c, db)
		})

	}
}

func setupTruexRoutes(baseGroup *gin.RouterGroup, redis database.RedisClientInterface, db *gorm.DB) {
	truexGroup := baseGroup.Group("/truex")
	truexGroup.Use()
	{
		truexGroup.GET("/packages", func(c *gin.Context) {
			handlers.GetPackagesList(c, redis, db)
		})

	}
}
