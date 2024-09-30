package utils

import "github.com/gofiber/fiber/v2"

func SendErrorResponse(c *fiber.Ctx, statusCode int, message string) error {
	return c.Status(statusCode).JSON(fiber.Map{
		"status": false,
		"error": fiber.Map{
			"statusCode": statusCode,
			"message":    message,
		},
	})
}

func SendSuccessResponse(c *fiber.Ctx, statusCode int, data interface{}) error {
	return c.Status(statusCode).JSON(fiber.Map{
		"status":     true,
		"statusCode": statusCode,
		"data":       data,
	})
}
