package handlers

import (
	"chat-api/config"
	"chat-api/middleware"
	"chat-api/models"
	"database/sql"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetChats(c *fiber.Ctx) error {
	rows, err := config.DB.Query(`
		SELECT chat_id, user_id, created_at, disease, text, name, height, weight,
		       gender, physical_condi, medical_histor, L, O, D, C, R, A, F, T
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
		err := rows.Scan(&chat.ChatID, &chat.UserID, &chat.CreatedAt,
			&chat.Disease, &chat.Text, &chat.Name, &chat.Height, &chat.Weight,
			&chat.Gender, &chat.PhysicalCondi, &chat.MedicalHistor,
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
	err = config.DB.QueryRow(`
		SELECT chat_id, user_id, created_at, disease, text, name, height, weight,
		       gender, physical_condi, medical_histor, L, O, D, C, R, A, F, T
		FROM chats WHERE chat_id = $1`, chatID).Scan(
		&chat.ChatID, &chat.UserID, &chat.CreatedAt,
		&chat.Disease, &chat.Text, &chat.Name, &chat.Height, &chat.Weight,
		&chat.Gender, &chat.PhysicalCondi, &chat.MedicalHistor,
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
	_, err := config.DB.Exec(`
		INSERT INTO chats (chat_id, user_id, created_at, disease, text, name, height, weight,
		                   gender, physical_condi, medical_histor, L, O, D, C, R, A, F, T)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)`,
		chatID, userID, time.Now(),
		input.Disease, input.Text, input.Name, input.Height, input.Weight,
		input.Gender, input.PhysicalCondi, input.MedicalHistor,
		input.L, input.O, input.D, input.C, input.R, input.A, input.F, input.T)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create chat",
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
	err = config.DB.QueryRow("SELECT user_id FROM chats WHERE chat_id = $1", chatID).Scan(&ownerID)
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

	if ownerID != userID {
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

	_, err = config.DB.Exec(`
		UPDATE chats SET disease = $1, text = $2, name = $3, height = $4, weight = $5,
		               gender = $6, physical_condi = $7, medical_histor = $8,
		               L = $9, O = $10, D = $11, C = $12, R = $13, A = $14, F = $15, T = $16
		WHERE chat_id = $17`,
		input.Disease, input.Text, input.Name, input.Height, input.Weight,
		input.Gender, input.PhysicalCondi, input.MedicalHistor,
		input.L, input.O, input.D, input.C, input.R, input.A, input.F, input.T, chatID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update chat",
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
	err = config.DB.QueryRow("SELECT user_id FROM chats WHERE chat_id = $1", chatID).Scan(&ownerID)
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

	if ownerID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You can only delete your own chats",
		})
	}

	_, err = config.DB.Exec("DELETE FROM chats WHERE chat_id = $1", chatID)
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
	rows, err := config.DB.Query(`
		SELECT chat_id, user_id, created_at, disease, text, name, height, weight,
		       gender, physical_condi, medical_histor, L, O, D, C, R, A, F, T
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
		err := rows.Scan(&chat.ChatID, &chat.UserID, &chat.CreatedAt,
			&chat.Disease, &chat.Text, &chat.Name, &chat.Height, &chat.Weight,
			&chat.Gender, &chat.PhysicalCondi, &chat.MedicalHistor,
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
