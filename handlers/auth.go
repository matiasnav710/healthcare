package handlers

import (
	"chat-api/database"
	"chat-api/middleware"
	"chat-api/models"
	"chat-api/utils"
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func SignUp(c *fiber.Ctx) error {
	var input models.UserSignUp
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	// Validate required fields
	if input.Email == "" || input.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email and password are required",
		})
	} else if len(input.Password) < 6 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Password must be at least 6 characters long",
		})
	}

	// Check if user already exists
	var existingEmail string
	err := database.DB.QueryRow("SELECT email FROM users WHERE email = $1", input.Email).Scan(existingEmail)
	if err != sql.ErrNoRows {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "User already exists " + err.Error(),
		})
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to hash password",
		})
	}

	// Create user
	var userID uuid.UUID
	err = database.DB.QueryRow(
		`INSERT INTO users (email, password) VALUES ($1, $2)
		 RETURNING user_id`,
		input.Email, hashedPassword).Scan(&userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user" + err.Error(),
		})
	}

	// Generate JWT
	token, err := middleware.GenerateJWTToken(userID, input.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User created successfully",
		"token":   token,
		"user": fiber.Map{
			"user_id": userID,
			"email":   input.Email,
		},
	})
}

func SignIn(c *fiber.Ctx) error {
	// Decode JWT
	var user models.User
	token, tokenEerr := middleware.DecodeJWTTokenFromHeader(c)
	if tokenEerr == nil {
		database.DB.QueryRow(`
		SELECT user_id, email
		FROM users WHERE user_id = $1 AND email = $2`, token.UserID, token.Email).Scan(
			&user.UserID, &user.Email)
		if user.UserID != uuid.Nil && user.Email != "" {
			return c.JSON(fiber.Map{
				"message": "Login successful",
				"token":   token,
				"user": fiber.Map{
					"user_id": user.UserID,
					"email":   user.Email,
					"name":    user.Name,
				},
			})
		}
	}

	var input models.UserSignIn
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	err := database.DB.QueryRow(`
		SELECT user_id, email, password, name 
		FROM users WHERE email = $1`, input.Email).Scan(
		&user.UserID, &user.Email, &user.Password, &user.Name)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	// Check password
	if !utils.CheckPasswordHash(input.Password, user.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "wrong password",
		})
	}

	// Generate JWT
	genToken, err := middleware.GenerateJWTToken(user.UserID, user.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Login successful",
		"token":   genToken,
		"user": fiber.Map{
			"user_id": user.UserID,
			"email":   user.Email,
			"name":    user.Name,
		},
	})
}
