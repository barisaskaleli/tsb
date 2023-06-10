package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"tsb/handler"
)

func SetRoutes(app *fiber.App) {
	app.Use(compress.New())

	api := app.Group("/api")

	// set method to post and retrieve 1 variable that is year (int) and path to /get-files
	api.Post("/get-files", handler.RetrieveFiles)
}
