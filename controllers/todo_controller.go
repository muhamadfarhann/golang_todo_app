package controllers

import (
	"net/http"
	"time"
	"todo_app/config"
	"todo_app/middleware"
	"todo_app/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Register user
func Register(c *gin.Context) {
	var input models.RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := models.User{
		Username: input.Username,
		Password: string(hashedPassword),
		Role:     input.Role,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"role":     user.Role,
		},
	})
}

// Login user
func Login(c *gin.Context) {
	var input models.LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := config.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	claims := &middleware.Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(middleware.JWTSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"role":     user.Role,
		},
	})
}

// Create Todo (User only)
func CreateTodo(c *gin.Context) {
	var input models.TodoInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userID")

	todo := models.Todo{
		UserID:      userID.(uint),
		Judul:       input.Judul,
		Deskripsi:   input.Deskripsi,
		Kategori:    input.Kategori,
		Priority:    input.Priority,
		IsCompleted: input.IsCompleted,
	}

	if err := config.DB.Create(&todo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create todo"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Todo created successfully",
		"todo":    todo,
	})
}

// Get All Todos (User: own todos, Admin: all todos)
func GetTodos(c *gin.Context) {
	var todos []models.Todo
	userID, _ := c.Get("userID")
	role, _ := c.Get("role")

	query := config.DB.Preload("User")

	// User can only see their own todos
	if role.(string) == "user" {
		query = query.Where("user_id = ?", userID.(uint))
	}
	// Admin can see all todos

	if err := query.Find(&todos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch todos"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"todos": todos})
}

// Get Todo by ID
func GetTodoByID(c *gin.Context) {
	id := c.Param("id")
	var todo models.Todo

	userID, _ := c.Get("userID")
	role, _ := c.Get("role")

	query := config.DB.Preload("User")

	if role.(string) == "user" {
		query = query.Where("user_id = ?", userID.(uint))
	}

	if err := query.First(&todo, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"todo": todo})
}

// Update Todo (User only, own todos)
func UpdateTodo(c *gin.Context) {
	id := c.Param("id")
	var todo models.Todo

	userID, _ := c.Get("userID")

	// Check if todo exists and belongs to user
	if err := config.DB.Where("id = ? AND user_id = ?", id, userID.(uint)).First(&todo).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found or unauthorized"})
		return
	}

	var input models.TodoInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := models.Todo{
		Judul:       input.Judul,
		Deskripsi:   input.Deskripsi,
		Kategori:    input.Kategori,
		Priority:    input.Priority,
		IsCompleted: input.IsCompleted,
	}

	if err := config.DB.Model(&todo).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update todo"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Todo updated successfully",
		"todo":    todo,
	})
}

// Delete Todo (User only, own todos)
func DeleteTodo(c *gin.Context) {
	id := c.Param("id")
	var todo models.Todo

	userID, _ := c.Get("userID")

	// Check if todo exists and belongs to user
	if err := config.DB.Where("id = ? AND user_id = ?", id, userID.(uint)).First(&todo).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found or unauthorized"})
		return
	}

	// // Soft Delete
	// if err := config.DB.Delete(&todo).Error; err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete todo"})
	// 	return
	// }

	// Permanent delete dengan Unscoped()
	if err := config.DB.Unscoped().Delete(&todo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete todo"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Todo deleted successfully"})
}
