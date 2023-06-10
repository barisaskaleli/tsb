package main

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"tsb/config"
	router "tsb/route"
)

func main() {
	fiberConfig := fiber.Config{
		AppName:     "[TSB API]",
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	}

	app := fiber.New(fiberConfig)

	config.DBConnect()

	router.SetRoutes(app)

	err := app.Listen(":3000")

	if err != nil {
		panic(err)
	}
}
