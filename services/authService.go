package services

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/Eberewill/emotechat/initializers"
	"github.com/Eberewill/emotechat/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	// Check if a user with the given email already exists
	existingUser := models.User{}
	if err := initializers.DB.Where("email = ?", data["email"]).First(&existingUser).Error; err == nil {
		// User with the same email already exists
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "User with this email already exists",
		})
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)

	user := models.User{
		Name:     data["name"],
		Email:    data["email"],
		Password: password,
	}

	initializers.DB.Create(&user)

	return c.JSON(user)
}

func Login(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	var user models.User
	initializers.DB.Where("email = ?", data["email"]).First(&user)

	if user.Id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "User not found",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data["password"])); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Incorrect password",
		})
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   strconv.Itoa(int(user.Id)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // 1 day
	})

	token, err := claims.SignedString([]byte(os.Getenv("SECRETEKEY")))
	if err != nil {
		fmt.Println("Token creation error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not log in",
		})
	}
	/*
		cookie := fiber.Cookie{
			Name:     "jwt",
			Value:    token,
			Expires:  time.Now().Add(time.Hour * 24),
			HTTPOnly: true,
		}

		c.Cookie(&cookie)
	*/

	c.Response().Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	return c.JSON(fiber.Map{
		"message": "Success",
		"user":    user,
		"token":   fmt.Sprintf("Bearer %s", token),
	})
}

func Logout(c *fiber.Ctx) error {

	// get the session from the authorization header
	sessionHeader := c.Get("Authorization")

	// ensure the session header is not empty and in the correct format
	if sessionHeader == "" || len(sessionHeader) < 8 || sessionHeader[:7] != "Bearer " {
		return c.JSON(fiber.Map{"error": "invalid session header"})
	}

	c.Response().Header.Del("Authorization")

	return c.JSON(fiber.Map{
		"message": "Logged out successfully",
	})
}

func Validate(c *fiber.Ctx) error {
	// Retrieve the userID from context locals
	userID := c.Locals("userID").(string)

	var user models.User
	initializers.DB.Where("id = ?", userID).First(&user)

	if user.Id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "User not found",
		})
	}

	// Now you have the userID of the current logged-in user
	return c.JSON(fiber.Map{
		"message": "Validated",
		"user":    user,
	})
}
