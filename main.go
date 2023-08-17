package main

import (
	"TESTGO/api"
	"TESTGO/pkg/database/mysql"
	"TESTGO/pkg/database/redis"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	file, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		logrus.SetOutput(file)
	} else {
		logrus.Info("Failed to log to file, using default stderr")
	}

	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetReportCaller(true)
}

func main() {

	db := mysql.Initialize()        // Initialize Database
	redis := redis.NewRedisClient() // Initialize redis
	defer redis.CloseRedis()        // close redis connection
	api.StartAPI(redis, db)         // StartApi รับค่า redis และ db เพื่อเริ่มต้นระบบ API

}
