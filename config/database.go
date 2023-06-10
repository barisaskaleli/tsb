package config

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"tsb/helper"
)

func DBConnect() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", helper.GetEnv("DB_USERNAME"), helper.GetEnv("DB_PASSWORD"), helper.GetEnv("DB_HOST"), helper.GetEnv("DB_PORT"), helper.GetEnv("DB_DATABASE"))

	Database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	Database.Set("gorm:table_options", "CHARSET=utf8mb4")

	fmt.Println("Connection Opened to Database")
}
