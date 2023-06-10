package handler

import (
	"github.com/gofiber/fiber/v2"
	TSBService "tsb/service"
)

func RetrieveFiles(ctx *fiber.Ctx) error {
	result := TSBService.GetData(ctx)

	return result
}
