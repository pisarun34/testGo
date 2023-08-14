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

	db := mysql.Initialize() // Initialize ฐานข้อมูล
	redis := redis.NewRedisClient()
	defer redis.CloseRedis()
	api.StartAPI(redis, db)

}
