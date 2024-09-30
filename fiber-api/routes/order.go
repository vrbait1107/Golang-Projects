package routes

import (
	"errors"
	"time"

	"example.com/package/database"
	"example.com/package/models"
	"example.com/package/utils"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
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
		return utils.SendErrorResponse(c, fiber.StatusNotFound, "User record not found")
	}

	var product models.Product

	if err := database.Database.Db.First(&product, order.ProductRefer).Error; err != nil {
		return utils.SendErrorResponse(c, fiber.StatusNotFound, "Product record not found")
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

func GetOrder(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid ID parameter")
	}

	var order models.Order

	if err := database.Database.Db.First(&order, id).Error; err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Internal server error")
	}

	var product models.Product
	var user models.User

	if err := database.Database.Db.Find(&user, order.UserRefer).Error; err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Internal server error")
	}

	if err := database.Database.Db.Find(&product, order.ProductRefer).Error; err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Internal server error")
	}

	return utils.SendSuccessResponse(c, fiber.StatusOK, CreateResponseOrder(order, CreateResponseUser(user), CreateResponseProduct(product)))
}

func DeleteOrder(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid ID parameter")
	}

	var order models.Order

	result := database.Database.Db.First(&order, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return utils.SendErrorResponse(c, fiber.StatusNotFound, "Order not found")
		}

		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Internal server error")
	}

	if err := database.Database.Db.Delete(&order).Error; err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Internal server error")
	}

	return utils.SendSuccessResponse(c, fiber.StatusOK, nil)

}
