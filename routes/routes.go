package routes

import (
	"chat-api/handlers"
	"chat-api/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	//default route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to Chat API")
	})
	// Auth routes (public)
	// auth routes don't require JWT token
	auth := app.Group("/auth")
	auth.Post("/signup", handlers.SignUp) //body required: email, password(6 length)
	auth.Post("/signin", handlers.SignIn) //body required: email, password(6 length)

	// Protected routes
	api := app.Group("/api")
	protected := api.Group("", middleware.SetJWtHeaderHandler()) //all below route require JWT token

	// User routes
	users := protected.Group("/users")
	users.Get("/", handlers.GetUsers)   // Get all users
	users.Get("/:id", handlers.GetUser) // Get user by ID
	// Create a new user (jwt must role admin) | body required: email, password(6 length)
	users.Post("/create", handlers.CreateUser)
	//jwt must role admin or have the same user ID as params | can update role to admin if and only if jwt role is admin
	users.Put("/:id", handlers.UpdateUser)
	// Delete a user by ID (jwt must role admin)
	users.Delete("/:id", handlers.DeleteUser)

	// Chat routes
	chats := protected.Group("/chats")
	chats.Get("/", handlers.GetChats)               // get all chats
	chats.Get("/getByChatID/:id", handlers.GetChat) // Get chat by ID
	chats.Post("/", handlers.CreateChat)            // Create a new chat
	// Update a chat by ID (jwt must role admin or have the same user ID as chat's user_id)
	chats.Put("/:id", handlers.UpdateChat)
	// Delete a chat by ID (jwt must role admin or have the same user ID as chat's user_id)
	chats.Delete("/:id", handlers.DeleteChat)
	chats.Get("/all_chat_id", handlers.GetUserChats) // Get user's all chats
}
