package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserID        uuid.UUID `json:"user_id" db:"user_id"`
	Email         string    `json:"email" db:"email"`
	Password      string    `json:"password,omitempty" db:"password"`
	Name          *string   `json:"name" db:"name"`
	Age           *int      `json:"age" db:"age"`
	Height        *float64  `json:"height" db:"height"`
	Weight        *float64  `json:"weight" db:"weight"`
	Gender        *string   `json:"gender" db:"gender"`
	PhysicalCondi *string   `json:"physical_condi" db:"physical_condi"`
	MedicalHistor *string   `json:"medical_histor" db:"medical_histor"`
	ProfileImage  *string   `json:"profile_image_" db:"profile_image_"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
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
	UserID        uuid.UUID `json:"user_id"`
	Email         string    `json:"email"`
	Name          *string   `json:"name"`
	Age           *int      `json:"age"`
	Height        *float64  `json:"height"`
	Weight        *float64  `json:"weight"`
	Gender        *string   `json:"gender"`
	PhysicalCondi *string   `json:"physical_condi"`
	MedicalHistor *string   `json:"medical_histor"`
	ProfileImage  *string   `json:"profile_image_"`
	CreatedAt     time.Time `json:"created_at"`
}
