package utils

import (
	"chat-api/models"
	"strconv"
)

func BuildUsersUpdateDynamicArray(data *models.UserInsertUpdate, role string) (string, []interface{}, int, error) {
	query := "UPDATE users SET "
	args := []interface{}{}
	argCount := 1

	if data.Email != "" {
		query += "email = $" + strconv.Itoa(argCount) + ", "
		args = append(args, data.Email)
		argCount++
	}
	if data.Password != "" {
		hashedPassword, err := HashPassword(data.Password)
		if err != nil {
			return "", nil, 0, err
		}
		query += "password = $" + strconv.Itoa(argCount) + ", "
		args = append(args, hashedPassword)
		argCount++
	}
	if role == "admin" && data.Role != "" {
		query += "role = $" + strconv.Itoa(argCount) + ", "
		args = append(args, data.Role)
		argCount++
	}
	if data.Name != nil {
		query += "name = $" + strconv.Itoa(argCount) + ", "
		args = append(args, *data.Name)
		argCount++
	}
	if data.Age != nil {
		query += "age = $" + strconv.Itoa(argCount) + ", "
		args = append(args, *data.Age)
		argCount++
	}
	if data.Height != nil {
		query += "height = $" + strconv.Itoa(argCount) + ", "
		args = append(args, *data.Height)
		argCount++
	}
	if data.Weight != nil {
		query += "weight = $" + strconv.Itoa(argCount) + ", "
		args = append(args, *data.Weight)
		argCount++
	}
	if data.Gender != nil {
		query += "gender = $" + strconv.Itoa(argCount) + ", "
		args = append(args, *data.Gender)
		argCount++
	}
	if data.PhysicalCondition != nil {
		query += "physical_condition = $" + strconv.Itoa(argCount) + ", "
		args = append(args, *data.PhysicalCondition)
		argCount++
	}
	if data.MedicalHistory != nil {
		query += "medical_history = $" + strconv.Itoa(argCount) + ", "
		args = append(args, *data.MedicalHistory)
		argCount++
	}
	if data.ProfileImageUrl != nil {
		query += "profile_image_url = $" + strconv.Itoa(argCount) + ", "
		args = append(args, *data.ProfileImageUrl)
		argCount++
	}

	// Remove trailing comma and space
	query = query[:len(query)-2]
	return query, args, argCount, nil
}
