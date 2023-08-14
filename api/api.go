package api

import (
	"TESTGO/pkg/database"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// StartAPI จะเริ่มต้นระบบ API โดยรับค่า db เพื่อเชื่อมต่อฐานข้อมูล
func StartAPI(redis database.RedisClientInterface, db *gorm.DB) {
	r := gin.Default()
	SetupRouter(r, redis, db)
	r.Run(":8080")
}
