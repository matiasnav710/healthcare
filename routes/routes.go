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
	users.Put("/:id", handlers.UpdateUser)
	users.Delete("/:id", handlers.DeleteUser)

	// Chat routes
	chats := protected.Group("/chats")
	chats.Get("/", handlers.GetChats)
	chats.Get("/:id", handlers.GetChat)
	chats.Post("/", handlers.CreateChat)
	chats.Put("/:id", handlers.UpdateChat)
	chats.Delete("/:id", handlers.DeleteChat)
	chats.Get("/my", handlers.GetUserChats) // Get current user's chats

	// User-Chat relation routes
	userChats := protected.Group("/user-chats")
	userChats.Get("/", handlers.GetUserChatsRelations)
	userChats.Post("/", handlers.CreateUserChatRelation)
	userChats.Get("/:user_id/:chat_id", handlers.GetUserChatRelation)
	userChats.Delete("/:user_id/:chat_id", handlers.DeleteUserChatRelation)
}
