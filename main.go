package main

import (
	"TESTGO/api"
	"TESTGO/pkg/database/mysql"
	"TESTGO/pkg/database/redis"
	"log"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {

	db := mysql.Initialize()        // Initialize Database
	redis := redis.NewRedisClient() // Initialize redis
	defer redis.CloseRedis()        // close redis connection
	api.StartAPI(redis, db)         // StartApi รับค่า redis และ db เพื่อเริ่มต้นระบบ API

}
