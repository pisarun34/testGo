package mysql

import (
	"TESTGO/pkg/database/models"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Initialize() *gorm.DB {
	// สร้าง connection string
	// โหลด config
	fmt.Println("InitializeMySQL")
	config := NewConfig()
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.User, config.Pass, config.Host, config.Name)
	//dsn := "user:password@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("err: ", err)
		panic("Failed to connect to database")
	}
	db.AutoMigrate(&models.Price{}, &models.Packages{}, &models.SubscriptionPackages{})
	// สร้างตารางโดยอัตโนมัติ
	return db
}
