package handlers

import (
	"chat-api/database"
	"chat-api/middleware"
	"chat-api/models"
	"database/sql"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetChats(c *fiber.Ctx) error {
	rows, err := database.DB.Query(`
		SELECT chat_id, user_id, created_at, updated_at, disease, text, name, age, height, weight,
		       blood_pressure, pulse, gender, physical_condition, medical_history,
			   "L", "O", "D", "C", "R", "A", "F", "T"
		FROM chats ORDER BY created_at DESC`)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch chats",
		})
	}
	defer rows.Close()

	var chats []models.Chat
	for rows.Next() {
		var chat models.Chat
		err := rows.Scan(&chat.ChatID, &chat.UserID, &chat.CreatedAt, &chat.UpdatedAt,
			&chat.Disease, &chat.Text, &chat.Name, &chat.Age, &chat.Height, &chat.Weight,
			&chat.BloodPressure, &chat.Pulse, &chat.Gender, &chat.PhysicalCondition, &chat.MedicalHistory,
			&chat.L, &chat.O, &chat.D, &chat.C, &chat.R, &chat.A, &chat.F, &chat.T)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to scan chat",
			})
		}
		chats = append(chats, chat)
	}

	return c.JSON(chats)
}

func GetChat(c *fiber.Ctx) error {
	chatID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid chat ID",
		})
	}

	var chat models.Chat
	err = database.DB.QueryRow(`
		SELECT chat_id, user_id, created_at, updated_at, disease, text, name, age, height, weight,
		       blood_pressure, pulse, gender, physical_condition, medical_history,
			   "L", "O", "D", "C", "R", "A", "F", "T"
		FROM chats WHERE chat_id = $1`, chatID).Scan(
		&chat.ChatID, &chat.UserID, &chat.CreatedAt, &chat.UpdatedAt,
		&chat.Disease, &chat.Text, &chat.Name, &chat.Age, &chat.Height, &chat.Weight,
		&chat.BloodPressure, &chat.Pulse, &chat.Gender, &chat.PhysicalCondition, &chat.MedicalHistory,
		&chat.L, &chat.O, &chat.D, &chat.C, &chat.R, &chat.A, &chat.F, &chat.T)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Chat not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch chat",
		})
	}

	return c.JSON(chat)
}

func CreateChat(c *fiber.Ctx) error {
	td, tokenErr := middleware.DecodeJWTToken(c)
	if tokenErr != nil {
		return tokenErr
	}

	userID := td.UserID
	var input models.ChatCreate

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	chatID := uuid.New()
	_, err := database.DB.Exec(`
		INSERT INTO chats (chat_id, user_id, created_at, updated_at, disease, text, name, age, height, weight,
		                   blood_pressure, pulse, gender, physical_condition, medical_history,
						   "L", "O", "D", "C", "R", "A", "F", "T")
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23)`,
		chatID, userID, time.Now(), time.Now(),
		input.Disease, input.Text, input.Name, input.Age, input.Height, input.Weight,
		input.BloodPressure, input.Pulse, input.Gender, input.PhysicalCondition, input.MedicalHistory,
		input.L, input.O, input.D, input.C, input.R, input.A, input.F, input.T)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create chat " + err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Chat created successfully",
		"chat_id": chatID,
	})
}

func UpdateChat(c *fiber.Ctx) error {
	td, err := middleware.DecodeJWTToken(c)
	if err != nil {
		return err
	}

	userID := td.UserID
	chatID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid chat ID",
		})
	}

	// Check if chat belongs to user
	var ownerID uuid.UUID
	err = database.DB.QueryRow("SELECT user_id FROM chats WHERE chat_id = $1", chatID).Scan(&ownerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Chat not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to verify chat ownership",
		})
	}

	if td.Role != "admin" && ownerID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You can only update your own chats",
		})
	}

	var input models.ChatCreate
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	_, err = database.DB.Exec(`
		UPDATE chats SET updated_at = $1, disease = $2, text = $3, name = $4, age = $5, height = $6, weight = $7,
		               blood_pressure = $8, pulse = $9, gender = $10, physical_condition = $11, medical_history = $12,
		               "L" = $13, "O" = $14, "D" = $15, "C" = $16, "R" = $17, "A" = $18, "F" = $19, "T" = $20
		WHERE chat_id = $21`,
		time.Now(), input.Disease, input.Text, input.Name, input.Age, input.Height, input.Weight,
		input.BloodPressure, input.Pulse, input.Gender, input.PhysicalCondition, input.MedicalHistory,
		input.L, input.O, input.D, input.C, input.R, input.A, input.F, input.T, chatID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update chat " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Chat updated successfully",
	})
}

func DeleteChat(c *fiber.Ctx) error {
	td, err := middleware.DecodeJWTToken(c)
	if err != nil {
		return err
	}

	userID := td.UserID
	chatID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid chat ID",
		})
	}

	// Check if chat belongs to user
	var ownerID uuid.UUID
	err = database.DB.QueryRow("SELECT user_id FROM chats WHERE chat_id = $1", chatID).Scan(&ownerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Chat not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to verify chat ownership",
		})
	}

	if td.Role != "admin" && ownerID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You can only delete your own chats",
		})
	}

	_, err = database.DB.Exec("DELETE FROM chats WHERE chat_id = $1", chatID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete chat",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Chat deleted successfully",
	})
}

func GetUserChats(c *fiber.Ctx) error {
	td, err := middleware.DecodeJWTToken(c)
	if err != nil {
		return err
	}

	userID := td.UserID
	rows, err := database.DB.Query(`
		SELECT chat_id, user_id, created_at, updated_at, disease, text, name, age, height, weight,
		       blood_pressure, pulse, gender, physical_condition, medical_history,
			   "L", "O", "D", "C", "R", "A", "F", "T"
		FROM chats WHERE user_id = $1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch user chats",
		})
	}
	defer rows.Close()

	var chats []models.Chat
	for rows.Next() {
		var chat models.Chat
		err := rows.Scan(&chat.ChatID, &chat.UserID, &chat.CreatedAt, &chat.UpdatedAt,
			&chat.Disease, &chat.Text, &chat.Name, &chat.Age, &chat.Height, &chat.Weight,
			&chat.BloodPressure, &chat.Pulse, &chat.Gender, &chat.PhysicalCondition, &chat.MedicalHistory,
			&chat.L, &chat.O, &chat.D, &chat.C, &chat.R, &chat.A, &chat.F, &chat.T)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to scan chat",
			})
		}
		chats = append(chats, chat)
	}

	return c.JSON(chats)
}
