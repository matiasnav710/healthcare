package middleware

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func SetJWtHeaderHandler() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			Key:    []byte(os.Getenv("JWT_SECRET_KEY")),
			JWTAlg: jwtware.HS256,
		},
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(http.StatusUnauthorized).SendString("Unauthorization Token.")
		},
	})
}

type TokenDetails struct {
	Token     *string   `json:"token"`
	UserID    uuid.UUID `json:"user_id"`
	Email     string    `json:"email"`
	Role      string    `json:"role"` // Optional role field
	ExpiresIn *int64    `json:"exp"`
}

func GenerateJWTToken(userID uuid.UUID, email string, role string) (*TokenDetails, error) {
	now := time.Now().UTC()

	td := &TokenDetails{
		ExpiresIn: new(int64),
		Token:     new(string),
	}

	*td.ExpiresIn = now.Add(time.Hour * 6).Unix()

	td.UserID = userID
	td.Email = email
	td.Role = role

	//ส่วนของ signature
	SigningKey := []byte(os.Getenv("JWT_SECRET_KEY"))

	//สร้าง payload
	atClaims := make(jwt.MapClaims)
	atClaims["user_id"] = userID.String()
	atClaims["email"] = email
	atClaims["role"] = role
	atClaims["exp"] = time.Now().Add(6 * time.Hour).Unix()
	// atClaims["exp"] = time.Now().Add(14 * 24 * time.Hour).Unix()
	atClaims["iat"] = time.Now().Unix()
	atClaims["nbf"] = time.Now().Unix()

	log.Println("New claims: ", atClaims)

	//สร้าง token
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims).SignedString(SigningKey)
	if err != nil {
		return nil, fmt.Errorf("create: sign token: %w", err)
	}

	*td.Token = token
	return td, nil
}

func DecodeJWTToken(ctx *fiber.Ctx) (*TokenDetails, error) {
	td := &TokenDetails{
		Token: new(string),
	}

	token, status := ctx.Locals("user").(*jwt.Token)
	if !status {
		return nil, ctx.Status(http.StatusUnauthorized).SendString("Unauthorization Token.")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ctx.Status(http.StatusUnauthorized).SendString("Unauthorization Token.")
	}

	for key, value := range claims {
		if key == "user_id" {
			userID, err := uuid.Parse(value.(string))
			if err != nil {
				return nil, ctx.Status(http.StatusUnauthorized).SendString("cannot parse user_id from token")
			}
			td.UserID = userID
		}
		if key == "email" {
			td.Email = value.(string)
		}
		if key == "role" {
			td.Role = value.(string)
		}
	}
	*td.Token = token.Raw
	return td, nil
}

func DecodeJWTTokenFromHeader(ctx *fiber.Ctx) (*TokenDetails, error) {
	authHeader := ctx.Get("Authorization")
	if authHeader == "" {
		return nil, fmt.Errorf("no authorization header")
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}

	td := &TokenDetails{
		Token: new(string),
	}
	for key, value := range claims {
		if key == "user_id" {
			userID, err := uuid.Parse(value.(string))
			if err != nil {
				return nil, ctx.Status(http.StatusUnauthorized).SendString("cannot parse user_id from token")
			}
			td.UserID = userID
		}
		if key == "email" {
			td.Email = value.(string)
		}
		if key == "role" {
			td.Role = value.(string)
		}
	}
	*td.Token = tokenStr
	return td, nil
}
