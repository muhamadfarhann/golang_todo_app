package routes

import (
	"todo_app/controllers"
	"todo_app/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	// Public routes
	auth := router.Group("/api/auth")
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
	}

	// Protected routes
	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		// User routes (create, update, delete)
		user := api.Group("/user")
		user.Use(middleware.RoleMiddleware("user"))
		{
			user.POST("/todos", controllers.CreateTodo)
			user.PUT("/todos/:id", controllers.UpdateTodo)
			user.DELETE("/todos/:id", controllers.DeleteTodo)
			user.GET("/todos", controllers.GetTodos)
			user.GET("/todos/:id", controllers.GetTodoByID)
		}

		// Admin routes (read only)
		admin := api.Group("/admin")
		admin.Use(middleware.RoleMiddleware("admin"))
		{
			admin.GET("/todos", controllers.GetTodos)
			admin.GET("/todos/:id", controllers.GetTodoByID)
		}
	}
}
