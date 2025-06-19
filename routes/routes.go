package routes

import (
	"chat-api/handlers"
	"chat-api/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	// Auth routes (public)
	auth := app.Group("/auth")
	auth.Post("/signup", handlers.SignUp)
	auth.Post("/signin", handlers.SignIn)

	// Protected routes
	api := app.Group("/api")
	protected := api.Group("", middleware.SetJWtHeaderHandler())

	// User routes
	users := protected.Group("/users")
	users.Get("/", handlers.GetUsers)
	users.Get("/:id", handlers.GetUser)
	users.Post("/create", handlers.CreateUser)
	users.Put("/:id", handlers.UpdateUser)
	users.Delete("/:id", handlers.DeleteUser)

	// Chat routes
	chats := protected.Group("/chats")
	chats.Get("/", handlers.GetChats)
	chats.Get("/getByChatID/:id", handlers.GetChat)
	chats.Post("/", handlers.CreateChat)
	chats.Put("/:id", handlers.UpdateChat)
	chats.Delete("/:id", handlers.DeleteChat)
	chats.Get("/all_chat_id", handlers.GetUserChats) // Get current user's chats
}
