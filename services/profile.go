package services

import (
	"fmt"
	"log"
	"math/rand"
	"net/smtp"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Eberewill/emotechat/initializers"
	"github.com/Eberewill/emotechat/models"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
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

func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := initializers.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByResetToken(token string) (*models.User, error) {
	var user models.User
	if err := initializers.DB.Where("reset_token = ? AND reset_token_expiry > ?", token, time.Now()).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func GenerateAndStoreResetToken(user *models.User) (string, error) {
	token := GenerateRandomToken()
	expiration := time.Now().Add(time.Hour * 24)

	user.ResetToken = token
	user.ResetTokenExpiry = expiration
	if err := initializers.DB.Save(user).Error; err != nil {
		return "", err
	}
	return token, nil
}

func SendResetEmail1(email, token string) error {
	from := os.Getenv("FROM_MAIL")
	password := os.Getenv("FROM_MAIL_PASSWORD")
	to := email

	resetURL := fmt.Sprintf("http://localhost:3000/reset-password?token=%s", token)
	message := fmt.Sprintf("Click the following link to reset your password: %s", resetURL)

	auth := smtp.PlainAuth("", from, password, "smtp.gmail.com")

	return smtp.SendMail("smtp.gmail.com:587", auth, from, []string{to}, []byte(message))
}

func SendResetEmail(email, token string) error {

	// Get sendgrid API Key
	sendgridAPIKey := os.Getenv("SENDGRID_API_KEY")
	host := os.Getenv("HOST")

	// Prepare the email components
	sender := mail.NewEmail(os.Getenv("APP_NAME"), os.Getenv("FROM_MAIL"))
	subject := "Password Reset Link"
	recipient := mail.NewEmail("User", email)

	resetURL := fmt.Sprintf("%s/reset-password?token=%s", host, token)
	message := fmt.Sprintf("Click the following link to reset your password: %s", resetURL)
	content := mail.NewContent("text/plain", message)

	// Initialize a new mail message
	m := mail.NewV3MailInit(sender, subject, recipient, content)

	// Initialize the SendGrid client
	client := sendgrid.NewSendClient(sendgridAPIKey)

	// Send the email
	_, err := client.Send(m)
	return err
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func UpdatePassword(user *models.User, newPassword string) error {
	return initializers.DB.Model(user).Updates(models.User{Password: []byte(newPassword), ResetToken: "", ResetTokenExpiry: time.Time{}}).Error
}

func GenerateRandomToken() string {
	const tokenLength = 32
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	token := make([]byte, tokenLength)
	for i := 0; i < tokenLength; i++ {
		token[i] = letters[rand.Intn(len(letters))]
	}
	return string(token)
}

// Route to initiate password reset
func RequestPasswordReset(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	var user models.User
	if err := initializers.DB.Where("email = ?", data["email"]).First(&user).Error; err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "User not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Database error"})
	}

	token, err := GenerateAndStoreResetToken(&user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to generate token"})
	}

	err = SendResetEmail(user.Email, token)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to send email"})
	}

	return c.JSON(fiber.Map{"message": "Password reset email sent"})
}

// Route to reset password
func ResetPassword(c *fiber.Ctx) error {
	token := c.FormValue("token")
	newPassword := c.FormValue("new_password")

	user, err := GetUserByResetToken(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid or expired token"})
	}

	hashedPassword, err := HashPassword(newPassword)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to hash password"})
	}

	err = UpdatePassword(user, hashedPassword)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to reset password"})
	}

	return c.JSON(fiber.Map{"message": "Password reset successfully"})
}
