package handlers

import (
	"chat-api/config"
	"chat-api/middleware"
	"chat-api/models"
	"chat-api/utils"
	"database/sql"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetUsers(c *fiber.Ctx) error {
	rows, err := config.DB.Query(`
		SELECT user_id, email, name, age, height, weight, gender, 
		       physical_condition , medical_history, profile_image_url
		FROM users`)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch users" + err.Error(),
		})
	}
	defer rows.Close()

	var users []models.UserResponse
	for rows.Next() {
		var user models.UserResponse
		err := rows.Scan(&user.UserID, &user.Email, &user.Name, &user.Age,
			&user.Height, &user.Weight, &user.Gender, &user.PhysicalCondition,
			&user.MedicalHistory, &user.ProfileImageUrl)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to scan user " + err.Error(),
			})
		}
		users = append(users, user)
	}

	return c.JSON(users)
}

func GetUser(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	var user models.UserResponse
	err = config.DB.QueryRow(`
		SELECT user_id, email, name, age, height, weight, gender,
		       physical_condition, medical_history, profile_image_url
		FROM users WHERE user_id = $1`, userID).Scan(
		&user.UserID, &user.Email, &user.Name, &user.Age,
		&user.Height, &user.Weight, &user.Gender, &user.PhysicalCondition,
		&user.MedicalHistory, &user.ProfileImageUrl)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch user " + err.Error(),
		})
	}

	return c.JSON(user)
}

func CreateUser(c *fiber.Ctx) error {
	// td, err := middleware.DecodeJWTToken(c)
	// if err != nil {
	// 	return err
	// }
	// if td.role != "admin" {
	// 	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
	// 		"error": "Only admins can create users",
	// 	})
	// }
	fmt.Println("I'm in insert user")
	var insertData models.UserInsertUpdate
	if err := c.BodyParser(&insertData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid updateData json",
		})
	}
	fmt.Println(insertData)
	var userID uuid.UUID
	err := config.DB.QueryRow(
		`INSERT INTO users
		(email, password, name, age, height, weight, gender, physical_condition, medical_history, profile_image_url)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING user_id`,
		insertData.Email, insertData.Password, insertData.Name,
		insertData.Age, insertData.Height, insertData.Weight, insertData.Gender,
		insertData.PhysicalCondition, insertData.MedicalHistory, insertData.ProfileImageUrl).Scan(&userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user " + err.Error(),
		})
	}
	token, err := middleware.GenerateJWTToken(userID, insertData.Email)
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
			"email":   insertData.Email,
		},
	})
}

func UpdateUser(c *fiber.Ctx) error {
	td, err := middleware.DecodeJWTToken(c)
	if err != nil {
		return err
	}

	userID := td.UserID
	paramID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	// Check if user is updating their own profile
	if userID != paramID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You can only update your own profile",
		})
	}

	updateData := models.UserInsertUpdate{}
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid updateData json",
		})
	}

	// Build dynamic query
	query, args, argCount, err := utils.BuildUsersUpdateDynamicArray(&updateData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to build update query: " + err.Error(),
		})
	}

	query += " WHERE user_id = $" + strconv.Itoa(argCount)
	args = append(args, userID)

	_, err = config.DB.Exec(query, args...)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update user " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "User updated successfully",
	})
}

func DeleteUser(c *fiber.Ctx) error {
	td, err := middleware.DecodeJWTToken(c)
	if err != nil {
		return err
	}

	userID := td.UserID
	paramID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	// Check if user is deleting their own account
	if userID != paramID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You can only delete your own account",
		})
	}

	_, err = config.DB.Exec("DELETE FROM users WHERE user_id = $1", userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete user",
		})
	}
	return c.JSON(fiber.Map{
		"message": "User deleted successfully",
	})
}
