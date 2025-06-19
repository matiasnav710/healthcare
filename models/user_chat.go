package models

import "github.com/google/uuid"

type UserChat struct {
	UserID uuid.UUID `json:"user_id" db:"user_id"`
	ChatID uuid.UUID `json:"chat_id" db:"chat_id"`
}

type UserChatCreate struct {
	UserID uuid.UUID `json:"user_id" validate:"required"`
	ChatID uuid.UUID `json:"chat_id" validate:"required"`
}