package models

import (
	"time"

	"github.com/google/uuid"
)

type Chat struct {
	ChatID        uuid.UUID  `json:"chat_id" db:"chat_id"`
	UserID        uuid.UUID  `json:"user_id" db:"user_id"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	Disease       *string    `json:"disease" db:"disease"`
	Text          *string    `json:"text" db:"text"`
	Name          *string    `json:"name" db:"name"`
	Height        *float64   `json:"height" db:"height"`
	Weight        *float64   `json:"weight" db:"weight"`
	Gender        *string    `json:"gender" db:"gender"`
	PhysicalCondi *string    `json:"physical_condi" db:"physical_condi"`
	MedicalHistor *string    `json:"medical_histor" db:"medical_histor"`
	L             *string    `json:"L" db:"L"`
	O             *string    `json:"O" db:"O"`
	D             *string    `json:"D" db:"D"`
	C             *string    `json:"C" db:"C"`
	R             *string    `json:"R" db:"R"`
	A             *string    `json:"A" db:"A"`
	F             *string    `json:"F" db:"F"`
	T             *string    `json:"T" db:"T"`
}

type ChatCreate struct {
	Disease       *string  `json:"disease"`
	Text          *string  `json:"text"`
	Name          *string  `json:"name"`
	Height        *float64 `json:"height"`
	Weight        *float64 `json:"weight"`
	Gender        *string  `json:"gender"`
	PhysicalCondi *string  `json:"physical_condi"`
	MedicalHistor *string  `json:"medical_histor"`
	L             *string  `json:"L"`
	O             *string  `json:"O"`
	D             *string  `json:"D"`
	C             *string  `json:"C"`
	R             *string  `json:"R"`
	A             *string  `json:"A"`
	F             *string  `json:"F"`
	T             *string  `json:"T"`
}