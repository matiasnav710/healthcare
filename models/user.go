package models

import (
	"github.com/google/uuid"
)

type User struct {
	UserID            uuid.UUID `json:"user_id" db:"user_id"`
	Email             string    `json:"email" db:"email"`
	Password          string    `json:"password,omitempty" db:"password"`
	Name              *string   `json:"name" db:"name"`
	Age               *int      `json:"age" db:"age"`
	Height            *float32  `json:"height" db:"height"`
	Weight            *float32  `json:"weight" db:"weight"`
	Gender            *string   `json:"gender" db:"gender"`
	PhysicalCondition *string   `json:"physical_condi" db:"physical_condition"`
	MedicalHistory    *string   `json:"medical_histor" db:"medical_history"`
	ProfileImageUrl   *string   `json:"profile_image_" db:"profile_image_url"`
}

type UserSignUp struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type UserSignIn struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserResponse struct {
	UserID            uuid.UUID `json:"user_id"`
	Email             string    `json:"email"`
	Name              *string   `json:"name"`
	Age               *int16    `json:"age"`
	Height            *float32  `json:"height"`
	Weight            *float32  `json:"weight"`
	Gender            *string   `json:"gender"`
	PhysicalCondition *string   `json:"physical_condition"`
	MedicalHistory    *string   `json:"medical_history"`
	ProfileImageUrl   *string   `json:"profile_image_url"`
}

type UserInsertUpdate struct {
	Email             string   `json:"email"`
	Password          string   `json:"password"`
	Name              *string  `json:"name"`
	Age               *int16   `json:"age"`
	Height            *float32 `json:"height"`
	Weight            *float32 `json:"weight"`
	Gender            *string  `json:"gender"`
	PhysicalCondition *string  `json:"physical_condition"`
	MedicalHistory    *string  `json:"medical_history"`
	ProfileImageUrl   *string  `json:"profile_image_url"`
}
