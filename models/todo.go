package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"unique;not null" json:"username"`
	Password  string         `gorm:"not null" json:"-"`
	Role      string         `gorm:"not null;default:'user'" json:"role"` // user or admin
	Todos     []Todo         `gorm:"foreignKey:UserID" json:"todos,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Todo struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UserID      uint           `gorm:"not null" json:"user_id"`
	User        User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Judul       string         `gorm:"not null" json:"judul"`
	Deskripsi   string         `gorm:"type:text" json:"deskripsi"`
	Kategori    string         `gorm:"not null" json:"kategori"` // General, Work, Personal, Shopping, Education, Health
	Priority    string         `gorm:"not null" json:"priority"` // high, medium, low
	IsCompleted bool           `gorm:"default:false" json:"is_completed"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required,oneof=user admin"`
}

type TodoInput struct {
	Judul       string `json:"judul" binding:"required"`
	Deskripsi   string `json:"deskripsi"`
	Kategori    string `json:"kategori" binding:"required,oneof=General Work Personal Shopping Education Health"`
	Priority    string `json:"priority" binding:"required,oneof=high medium low"`
	IsCompleted bool   `json:"is_completed"`
}
