package response

import "github.com/gofiber/fiber/v2"

type APIResponse interface {
	ResponseSuccess(c *fiber.Ctx, data interface{})
}

func ResponseSuccess(c *fiber.Ctx, data interface{}) error {
	return Response(c, fiber.StatusOK, fiber.Map{
		"status": "ok",
		"data":   data,
	})
}

func ResponseNotFound(c *fiber.Ctx, data interface{}) error {
	return Response(c, fiber.StatusNotFound, fiber.Map{
		"status": fiber.ErrNotFound.Message,
		"data":   data,
	})
}

func ResponseError(c *fiber.Ctx, data interface{}) error {
	return Response(c, fiber.StatusInternalServerError, fiber.Map{
		"status": fiber.ErrInternalServerError.Message,
		"data":   data,
	})
}

func ResponseBadRequest(c *fiber.Ctx, data interface{}) error {
	return Response(c, fiber.StatusBadRequest, fiber.Map{
		"status": fiber.ErrBadRequest.Message,
		"data":   data,
	})
}

func Response(c *fiber.Ctx, status int, data interface{}) error {
	return c.Status(status).JSON(data)
}
