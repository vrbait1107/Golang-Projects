package routes

import (
	"time"

	"example.com/package/database"
	"example.com/package/models"
	"example.com/package/utils"
	"github.com/gofiber/fiber/v2"
)

type OrderDTO struct {
	ID        uint       `json:"id"`
	User      UserDTO    `json:"user"`
	Product   ProductDTO `json:"product"`
	CreatedAt time.Time  `json:"order_date"`
}

func CreateResponseOrder(orderModel models.Order, user UserDTO, product ProductDTO) OrderDTO {
	return OrderDTO{ID: orderModel.ID, User: user, Product: product, CreatedAt: orderModel.CreatedAt}
}

func CreateOrder(c *fiber.Ctx) error {

	var order models.Order

	if err := c.BodyParser(&order); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid input data")
	}

	var user models.User

	if err := database.Database.Db.First(&user, order.UserRefer).Error; err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "User record not found")
	}

	var product models.Product

	if err := database.Database.Db.First(&product, order.ProductRefer).Error; err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Product record not found")
	}

	database.Database.Db.Create(&order)

	resposeUser := CreateResponseUser(user)
	responseProduct := CreateResponseProduct(product)
	responseOrder := CreateResponseOrder(order, resposeUser, responseProduct)

	return utils.SendSuccessResponse(c, fiber.StatusOK, responseOrder)

}

func GetOrders(c *fiber.Ctx) error {
	orders := []models.Order{}

	responseOrders := []OrderDTO{}

	database.Database.Db.Find(&orders)

	for _, order := range orders {
		var user models.User
		var product models.Product
		database.Database.Db.Find(&user, order.UserRefer)
		database.Database.Db.Find(&product, order.ProductRefer)
		responseOrder := CreateResponseOrder(order, CreateResponseUser(user), CreateResponseProduct(product))
		responseOrders = append(responseOrders, responseOrder)
	}

	return utils.SendSuccessResponse(c, fiber.StatusOK, responseOrders)
}
