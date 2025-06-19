package handlers

import (
	"chat-api/config"
	"chat-api/middleware"
	"chat-api/models"
	"chat-api/utils"
	"database/sql"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetUsers(c *fiber.Ctx) error {
	rows, err := config.DB.Query(`
		SELECT user_id, email, name, age, height, weight, gender, 
		       physical_condi, medical_histor, profile_image_, created_at
		FROM users ORDER BY created_at DESC`)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch users",
		})
	}
	defer rows.Close()

	var users []models.UserResponse
	for rows.Next() {
		var user models.UserResponse
		err := rows.Scan(&user.UserID, &user.Email, &user.Name, &user.Age,
			&user.Height, &user.Weight, &user.Gender, &user.PhysicalCondi,
			&user.MedicalHistor, &user.ProfileImage, &user.CreatedAt)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to scan user",
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
		       physical_condi, medical_histor, profile_image_, created_at
		FROM users WHERE user_id = $1`, userID).Scan(
		&user.UserID, &user.Email, &user.Name, &user.Age,
		&user.Height, &user.Weight, &user.Gender, &user.PhysicalCondi,
		&user.MedicalHistor, &user.ProfileImage, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch user",
		})
	}

	return c.JSON(user)
}

func UpdateUser(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
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

	var input struct {
		Name          *string  `json:"name"`
		Age           *int     `json:"age"`
		Height        *float64 `json:"height"`
		Weight        *float64 `json:"weight"`
		Gender        *string  `json:"gender"`
		PhysicalCondi *string  `json:"physical_condi"`
		MedicalHistor *string  `json:"medical_histor"`
		ProfileImage  *string  `json:"profile_image_"`
		Password      *string  `json:"password"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	// Build dynamic query
	query := "UPDATE users SET "
	args := []interface{}{}
	argCount := 1

	if input.Name != nil {
		query += "name = $" + strconv.Itoa(argCount) + ", "
		args = append(args, *input.Name)
		argCount++
	}
	if input.Age != nil {
		query += "age = $" + strconv.Itoa(argCount) + ", "
		args = append(args, *input.Age)
		argCount++
	}
	if input.Height != nil {
		query += "height = $" + strconv.Itoa(argCount) + ", "
		args = append(args, *input.Height)
		argCount++
	}
	if input.Weight != nil {
		query += "weight = $" + strconv.Itoa(argCount) + ", "
		args = append(args, *input.Weight)
		argCount++
	}
	if input.Gender != nil {
		query += "gender = $" + strconv.Itoa(argCount) + ", "
		args = append(args, *input.Gender)
		argCount++
	}
	if input.PhysicalCondi != nil {
		query += "physical_condi = $" + strconv.Itoa(argCount) + ", "
		args = append(args, *input.PhysicalCondi)
		argCount++
	}
	if input.MedicalHistor != nil {
		query += "medical_histor = $" + strconv.Itoa(argCount) + ", "
		args = append(args, *input.MedicalHistor)
		argCount++
	}
	if input.ProfileImage != nil {
		query += "profile_image_ = $" + strconv.Itoa(argCount) + ", "
		args = append(args, *input.ProfileImage)
		argCount++
	}
	if input.Password != nil {
		hashedPassword, err := utils.HashPassword(*input.Password)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to hash password",
			})
		}
		query += "password = $" + strconv.Itoa(argCount) + ", "
		args = append(args, hashedPassword)
		argCount++
	}

	// Remove trailing comma and space
	query = query[:len(query)-2]
	query += " WHERE user_id = $" + strconv.Itoa(argCount)
	args = append(args, userID)

	_, err = config.DB.Exec(query, args...)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update user",
		})
	}

	return c.JSON(fiber.Map{
		"message": "User updated successfully",
	})
}

func DeleteUser(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
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
