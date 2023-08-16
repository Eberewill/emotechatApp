package services

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/Eberewill/emotechat/initializers"
	"github.com/Eberewill/emotechat/models"
	"github.com/gofiber/fiber/v2"
)

// Create a struct to represent the paginated conversation response
type PaginatedConversations struct {
	Users         []models.User         `json:"users"`
	Conversations []models.Conversation `json:"conversations"`
	TotalPages    int                   `json:"total_pages"`
	CurrentPage   int                   `json:"current_page"`
	PageSize      int                   `json:"page_size"`
}

func FetchConversationsWithUsers(userID, partnerID uint, page, pageSize int) (*PaginatedConversations, error) {
	var result PaginatedConversations
	offset := (page - 1) * pageSize

	// Fetch involved users
	var users []models.User
	if err := initializers.DB.Where("id IN (?)", []uint{userID, partnerID}).Find(&users).Error; err != nil {
		return nil, err
	}

	var totalCount int64
	if err := initializers.DB.Model(&models.Conversation{}).
		Where("(receiver_id = ? AND sender_id = ?) OR (receiver_id = ? AND sender_id = ?)", userID, partnerID, partnerID, userID).
		Count(&totalCount).Error; err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))

	if err := initializers.DB.
		Where("(receiver_id = ? AND sender_id = ?) OR (receiver_id = ? AND sender_id = ?)", userID, partnerID, partnerID, userID).
		Order("timestamp ASC"). //  to ascending order here
		Offset(offset).
		Limit(pageSize).
		Preload("Receiver").
		Preload("Sender").
		Find(&result.Conversations).Error; err != nil {
		return nil, err
	}

	result.Users = users
	result.TotalPages = totalPages
	result.CurrentPage = page
	result.PageSize = pageSize

	return &result, nil
}

func SaveMessage(senderID, receiverID uint, content string) error {
	// Check if both sender and receiver exist in the database
	var sender, receiver models.User
	if err := initializers.DB.First(&sender, senderID).Error; err != nil {
		fmt.Println(err)
		return err
	}
	if err := initializers.DB.First(&receiver, receiverID).Error; err != nil {
		fmt.Println(err)
		return err
	}

	// Create and save the message
	message := models.Conversation{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Content:    content,
		Timestamp:  time.Now(),
	}
	if err := initializers.DB.Create(&message).Error; err != nil {
		fmt.Println(err)
		return err

	}
	return nil
}

func FetchConversations(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	userIDj, _ := strconv.Atoi(userID)

	// Parse query parameters for pagination
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.Query("page_size", "10"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	// Simulate partner user (replace with the partner's user ID from query params)
	partnerID, err := strconv.Atoi(c.Query("partner_id", "2"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid partner ID"})
	}

	// Fetch conversations with pagination
	paginatedConversations, err := FetchConversationsWithUsers(uint(userIDj), uint(partnerID), page, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch conversations"})
	}

	return c.JSON(paginatedConversations)
}

/*
func TestSaveMessage(c *fiber.Ctx) error {
	// Create a struct to represent the request payload
	type SaveMessageRequest struct {
		PartnerID uint `json:"partner_id"`
		//loggedInUserID uint   `json:"user_id"`
		MessageContent string `json:"message_content"`
	}

	// Parse the request payload
	var req SaveMessageRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	// Simulate logged-in user (replace with your actual authentication logic)

	// Save the message to the database
	err := SaveMessage(14, req.PartnerID, req.MessageContent)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save message"})
	}

	return c.JSON(fiber.Map{"message": "Message saved successfully"})
}

func TestFetchConversations(c *fiber.Ctx) error {
	// Simulate logged-in user (replace with your actual authentication logic)
	loggedInUserID := uint(14)

	// Parse the request payload
	var req struct {
		PartnerID uint `json:"partner_id"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	page := 1
	pageSize := 10

	// Fetch conversations
	conversations, err := FetchConversationsBetweenUsers(loggedInUserID, req.PartnerID, page, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch conversations"})
	}

	return c.JSON(conversations)
}
*/
