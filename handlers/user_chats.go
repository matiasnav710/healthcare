package handlers

import (
	"chat-api/config"
	"chat-api/middleware"
	"chat-api/models"
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetUserChatsRelations(c *fiber.Ctx) error {
	rows, err := config.DB.Query("SELECT user_id, chat_id FROM users_chats")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch user-chat relations",
		})
	}
	defer rows.Close()

	var userChats []models.UserChat
	for rows.Next() {
		var uc models.UserChat
		err := rows.Scan(&uc.UserID, &uc.ChatID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to scan user-chat relation",
			})
		}
		userChats = append(userChats, uc)
	}

	return c.JSON(userChats)
}

func CreateUserChatRelation(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	var input models.UserChatCreate

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	// Verify that the user is associating their own account or is admin
	if input.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You can only create relations for your own account",
		})
	}

	// Check if user exists
	var exists bool
	err := config.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE user_id = $1)", input.UserID).Scan(&exists)
	if err != nil || !exists {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User does not exist",
		})
	}

	// Check if chat exists
	err = config.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM chats WHERE chat_id = $1)", input.ChatID).Scan(&exists)
	if err != nil || !exists {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Chat does not exist",
		})
	}

	// Check if relation already exists
	err = config.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users_chats WHERE user_id = $1 AND chat_id = $2)", 
		input.UserID, input.ChatID).Scan(&exists)
	if err == nil && exists {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "User-chat relation already exists",
		})
	}

	_, err = config.DB.Exec("INSERT INTO users_chats (user_id, chat_id) VALUES ($1, $2)", 
		input.UserID, input.ChatID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user-chat relation",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User-chat relation created successfully",
	})
}

func DeleteUserChatRelation(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	relationUserID, err := uuid.Parse(c.Params("user_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	chatID, err := uuid.Parse(c.Params("chat_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid chat ID",
		})
	}

	// Check if user is deleting their own relation
	if relationUserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You can only delete your own user-chat relations",
		})
	}

	result, err := config.DB.Exec("DELETE FROM users_chats WHERE user_id = $1 AND chat_id = $2", 
		relationUserID, chatID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete user-chat relation",
		})
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get affected rows",
		})
	}

	if rowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User-chat relation not found",
		})
	}

	return c.JSON(fiber.Map{
		"message": "User-chat relation deleted successfully",
	})
}

func GetUserChatRelation(c *fiber.Ctx) error {
	relationUserID, err := uuid.Parse(c.Params("user_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	chatID, err := uuid.Parse(c.Params("chat_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid chat ID",
		})
	}

	var userChat models.UserChat
	err = config.DB.QueryRow("SELECT user_id, chat_id FROM users_chats WHERE user_id = $1 AND chat_id = $2", 
		relationUserID, chatID).Scan(&userChat.UserID, &userChat.ChatID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User-chat relation not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch user-chat relation",
		})
	}

	return c.JSON(userChat)
}