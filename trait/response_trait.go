package response_message_trait

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

const SUCCESS = "success"
const ERROR = "error"
const INVALID_REQUEST = "INVALID_REQUEST"
const UNPROCESSABLE_ENTITY = "UNPROCESSABLE_ENTITY"
const NOT_FOUND = "NOT_FOUND"

type Structure struct {
	Status  string            `json:"status"`
	Code    int               `json:"code"`
	Message interface{}       `json:"message"`
	Data    interface{}       `json:"data"`
	Errors  map[string]string `json:"errors"`
}

func Response(c *fiber.Ctx, r Structure) error {
	c.Status(r.Code)
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
	res, err := json.Marshal(r)

	if err != nil {
		return nil
	}

	return c.SendString(string(res))
}

func Success(c *fiber.Ctx, data interface{}) error {
	return Response(c, Structure{
		Status:  SUCCESS,
		Code:    http.StatusOK,
		Message: data,
		Data:    nil,
	})
}

func UnprocessableEntity(c *fiber.Ctx, errors map[string]string) error {
	return Response(c, Structure{
		Status:  ERROR,
		Code:    http.StatusUnprocessableEntity,
		Message: UNPROCESSABLE_ENTITY,
		Errors:  errors,
	})
}

func BadRequest(c *fiber.Ctx, errors map[string]string) error {
	return Response(c, Structure{
		Status:  ERROR,
		Code:    http.StatusBadRequest,
		Message: INVALID_REQUEST,
		Errors:  errors,
	})
}

func NotFound(c *fiber.Ctx, errors map[string]string) error {
	return Response(c, Structure{
		Status:  ERROR,
		Code:    http.StatusNotFound,
		Message: NOT_FOUND,
		Errors:  errors,
	})
}
