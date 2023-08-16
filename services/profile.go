package services

import (
	"log"
	"strconv"

	"github.com/Eberewill/emotechat/initializers"
	"github.com/Eberewill/emotechat/models"
	"github.com/gofiber/fiber/v2"
)

type Pagination struct {
	TotalPages  int `json:"total_pages"`
	CurrentPage int `json:"current_page"`
	PageSize    int `json:"page_size"`
}

func GetAllUsers(c *fiber.Ctx) error {
	db := initializers.DB

	// Pagination parameters
	pageSize := 40                           // Number of users per page
	page, _ := strconv.Atoi(c.Query("page")) // Convert page query parameter to integer
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * pageSize

	var users []models.User
	if err := db.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		log.Printf("Error fetching users: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "An error occurred while fetching users",
		})
	}

	var totalCount int64
	if err := db.Model(&models.User{}).Count(&totalCount).Error; err != nil {
		log.Printf("Error counting total users: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "An error occurred while counting total users",
		})
	}

	totalPages := int(totalCount) / pageSize
	if totalCount%int64(pageSize) > 0 {
		totalPages++
	}

	pagination := Pagination{
		TotalPages:  totalPages,
		CurrentPage: page,
		PageSize:    pageSize,
	}

	return c.JSON(fiber.Map{
		"users":      users,
		"pagination": pagination,
	})
}
