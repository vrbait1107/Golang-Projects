package main

import (
	"example.com/package/database"
	"example.com/package/routes"
	"example.com/package/utils"
	"github.com/gofiber/fiber/v2"
)

func WelcomeHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": true, "message": "Health OK"})
}

func main() {
	database.ConnectDb()
	app := fiber.New()
	app.Get("/", WelcomeHandler)
	SetupRoutes(app)
	app.Listen(":3000")
}

func SetupRoutes(app *fiber.App) {
	app.Get("/api", WelcomeHandler)
	app.Get("api/get-users", routes.GetUsers)
	app.Get("api/get-user/:id", routes.GetUser)
	app.Post("api/create-user", routes.CreateUser)
	app.Put("api/update-user/:id", routes.UpdateUsers)
	app.Delete("api/delete-user/:id", routes.DeleteUser)

	/* Product Routes */

	app.Get("api/get-products", routes.GetProducts)
	app.Get("api/get-product/:id", routes.GetProduct)
	app.Post("api/create-product", routes.CreateProduct)
	app.Put("api/update-product/:id", routes.UpdateProduct)
	app.Delete("api/delete-product/:id", routes.DeleteProduct)

	/* Order Routes */

	app.Get("api/get-orders", routes.GetOrders)
	app.Get("api/get-order/:id", routes.GetOrder)
	app.Post("api/create-order", routes.CreateOrder)
	app.Delete("api/delete-order/:id", routes.DeleteOrder)

	app.Use(func(c *fiber.Ctx) error {
		return utils.SendErrorResponse(c, fiber.StatusNotFound, "Route not found")
	})

}
