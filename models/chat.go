package models

import (
	"time"

	"github.com/google/uuid"
)

type Chat struct {
	ChatID            uuid.UUID `json:"chat_id" db:"chat_id"`
	UserID            uuid.UUID `json:"user_id" db:"user_id"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
	Disease           *string   `json:"disease" db:"disease"`
	Text              *string   `json:"text" db:"text"`
	Name              *string   `json:"name" db:"name"`
	Age               *int16    `json:"age" db:"age"`
	Height            *float32  `json:"height" db:"height"`
	Weight            *float32  `json:"weight" db:"weight"`
	BloodPressure     *string   `json:"blood_pressure" db:"blood_pressure"`
	Pulse             *int16    `json:"pulse" db:"pulse"`
	Gender            *string   `json:"gender" db:"gender"`
	PhysicalCondition *string   `json:"physical_condition" db:"physical_condition"`
	MedicalHistory    *string   `json:"medical_history" db:"medical_history"`
	L                 *string   `json:"L" db:"L"`
	O                 *string   `json:"O" db:"O"`
	D                 *string   `json:"D" db:"D"`
	C                 *string   `json:"C" db:"C"`
	R                 *string   `json:"R" db:"R"`
	A                 *string   `json:"A" db:"A"`
	F                 *string   `json:"F" db:"F"`
	T                 *string   `json:"T" db:"T"`
}

type ChatCreate struct {
	Disease           *string  `json:"disease"`
	Text              *string  `json:"text"`
	Name              *string  `json:"name"`
	Age               *int16   `json:"age" db:"age"`
	Height            *float32 `json:"height"`
	Weight            *float32 `json:"weight"`
	BloodPressure     *string  `json:"blood_pressure"`
	Pulse             *int16   `json:"pulse"`
	Gender            *string  `json:"gender"`
	PhysicalCondition *string  `json:"physical_condition"`
	MedicalHistory    *string  `json:"medical_history"`
	L                 *string  `json:"L"`
	O                 *string  `json:"O"`
	D                 *string  `json:"D"`
	C                 *string  `json:"C"`
	R                 *string  `json:"R"`
	A                 *string  `json:"A"`
	F                 *string  `json:"F"`
	T                 *string  `json:"T"`
}
