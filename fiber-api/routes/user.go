package routes

import (
	"errors"

	"example.com/package/database"
	"example.com/package/models"
	"example.com/package/utils"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UserDTO struct {
	ID        uint   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func CreateResponseUser(userModel models.User) UserDTO {
	return UserDTO{ID: userModel.ID, FirstName: userModel.FirstName, LastName: userModel.LastName}
}

type UpdateUser struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func CreateUser(c *fiber.Ctx) error {

	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid input data")
	}

	database.Database.Db.Create(&user)
	responeUser := CreateResponseUser(user)

	return utils.SendSuccessResponse(c, fiber.StatusOK, responeUser)

}

func GetUsers(c *fiber.Ctx) error {

	users := []models.User{}

	database.Database.Db.Find(&users)

	responseUsers := []UserDTO{}

	for _, user := range users {
		responseUser := CreateResponseUser(user)
		responseUsers = append(responseUsers, responseUser)
	}

	return utils.SendSuccessResponse(c, fiber.StatusOK, responseUsers)

}

func GetUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid ID parameter")
	}

	var user models.User

	result := database.Database.Db.First(&user, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return utils.SendErrorResponse(c, fiber.StatusNotFound, "User not found")
		}

		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Internal server error")
	}

	responseUser := CreateResponseUser(user)

	return utils.SendSuccessResponse(c, fiber.StatusOK, responseUser)
}

func UpdateUsers(c *fiber.Ctx) error {

	id, err := c.ParamsInt("id")

	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid ID parameter")
	}

	var user models.User

	result := database.Database.Db.First(&user, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return utils.SendErrorResponse(c, fiber.StatusNotFound, "User not found")
		}

		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Internal server error")
	}

	var updateData UpdateUser

	if err := c.BodyParser(&updateData); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid input data")
	}

	user.FirstName = updateData.FirstName
	user.LastName = updateData.LastName

	updateResponse := database.Database.Db.Save(&user)

	if updateResponse.Error != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Internal server error")
	}

	responseUser := CreateResponseUser(user)

	return utils.SendSuccessResponse(c, fiber.StatusOK, responseUser)

}

func DeleteUser(c *fiber.Ctx) error {

	id, err := c.ParamsInt("id")

	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Bad Request")
	}

	var user models.User

	result := database.Database.Db.First(&user, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return utils.SendErrorResponse(c, fiber.StatusNotFound, "User Not Found")
		}

		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Internal Server Error")
	}

	if err := database.Database.Db.Delete(&user).Error; err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Internal Server Error")
	}

	return utils.SendSuccessResponse(c, fiber.StatusOK, nil)

}
