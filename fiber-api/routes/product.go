package routes

import (
	"errors"

	"example.com/package/database"
	"example.com/package/models"
	"example.com/package/utils"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ProductDTO struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	SerialNumber string `json:"serial_number"`
}

type UpdateProductRequest struct {
	Name         string `json:"name"`
	SerialNumber string `json:"serial_number"`
}

func CreateResponseProduct(productModel models.Product) ProductDTO {
	return ProductDTO{
		ID:           productModel.ID,
		Name:         productModel.Name,
		SerialNumber: productModel.SerialNumber,
	}
}

func CreateProduct(c *fiber.Ctx) error {
	var product models.Product

	err := c.BodyParser(&product)

	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid input data")
	}

	product.ID = 0

	result := database.Database.Db.Create(&product)

	if result.Error != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Error creating product. Please try again later.")
	}

	return utils.SendSuccessResponse(c, fiber.StatusOK, product)
}

func GetProducts(c *fiber.Ctx) error {

	products := []models.Product{}

	database.Database.Db.Find(&products)

	responseProducts := []ProductDTO{}

	for _, product := range products {
		responseProduct := CreateResponseProduct(product)
		responseProducts = append(responseProducts, responseProduct)
	}

	return utils.SendSuccessResponse(c, fiber.StatusOK, responseProducts)

}

func GetProduct(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	if err != nil {
		utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid ID Parameter")
	}

	product := models.Product{}

	result := database.Database.Db.First(&product, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return utils.SendErrorResponse(c, fiber.StatusNotFound, "Product not found")
		}

		return utils.SendErrorResponse(c, fiber.StatusNotFound, "Internal server error")
	}

	responseProduct := CreateResponseProduct(product)

	return utils.SendSuccessResponse(c, fiber.StatusOK, responseProduct)

}

func UpdateProduct(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid ID parameter")
	}

	var product models.Product

	result := database.Database.Db.First(&product, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return utils.SendErrorResponse(c, fiber.StatusNotFound, "Product not found")
		}

		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Internal Server Error")
	}

	var updateProductRequest UpdateProductRequest

	if err := c.BodyParser(&updateProductRequest); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid input data")
	}

	product.Name = updateProductRequest.Name
	product.SerialNumber = updateProductRequest.SerialNumber

	updateProductResponse := database.Database.Db.Save(&product)

	if updateProductResponse.Error != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Internal server error")
	}

	resposneProduct := CreateResponseProduct(product)

	return utils.SendSuccessResponse(c, fiber.StatusOK, resposneProduct)

}

func DeleteProduct(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid ID parameter")
	}

	var product models.Product

	result := database.Database.Db.First(&product, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return utils.SendErrorResponse(c, fiber.StatusNotFound, "Product not found")
		}

		return utils.SendErrorResponse(c, fiber.StatusNotFound, "Internal server error")
	}

	if err := database.Database.Db.Delete(&product).Error; err != nil {
		return utils.SendErrorResponse(c, fiber.StatusNotFound, "Error deleting product, please try again later")
	}

	return utils.SendSuccessResponse(c, fiber.StatusOK, nil)
}
