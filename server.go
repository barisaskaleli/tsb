package main

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"tsb/config"
	"tsb/models"
	router "tsb/route"
)

func main() {
	fiberConfig := fiber.Config{
		AppName:     "[TSB API]",
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	}

	app := fiber.New(fiberConfig)

	db := config.DBConnect()

	err := db.AutoMigrate(
		&models.Brand{},
		&models.Model{},
		&models.CascoValue{},
	)
	if err != nil {
		fmt.Println("Error migrating database")
		panic(err.Error())
	}

	router.SetRoutes(app)

	err = app.Listen(":3000")

	if err != nil {
		panic(err)
	}
}
