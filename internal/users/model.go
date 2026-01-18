package users

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `gorm:"type:varchar(255);not null" json:"name"`
	Email    string `gorm:"type:varchar(255);not null;uniqueIndex" json:"email"`
	Password string `gorm:"type:varchar(255);not null" json:"-"` // "-" means don't include in JSON
	Token    string `gorm:"type:varchar(255)" json:"token,omitempty"`
}

func (User) TableName() string {
	return "users"
}
type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=3,max=100"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UpdateUserRequest struct {
	Name  string `json:"name" binding:"omitempty,min=3,max=100"`
	Email string `json:"email" binding:"omitempty,email"`
}

type UserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Token string `json:"token,omitempty"`
}
