package users

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
    ID        uint           `gorm:"column:user_id;primaryKey" json:"id"`
    Name      string         `gorm:"type:varchar(255);not null" json:"name"`
    Email     string         `gorm:"type:varchar(255);not null;uniqueIndex" json:"email"`
    Password  string         `gorm:"type:varchar(255);not null" json:"-"`
    Token     string         `gorm:"type:varchar(255)" json:"-"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (User) TableName() string {
	return "users"
}


type CreateUserRequest struct {
	ID       uint   `json:"id"`
	Name     string `json:"name" binding:"required,min=3,max=100"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=100"`
}


type UpdateUserRequest struct {
	Name  string `json:"name" binding:"omitempty,min=3,max=100"`
	Email string `json:"email" binding:"omitempty,email"`
}

type UserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
}

type CreateUserResonse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type GetUsersResponse struct {
	Data []User `json:"data"`
	Page int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
	TotalPages int `json:"total_pages"`
}