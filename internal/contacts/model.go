package contacts

import (
	"time"

	"github.com/DioSaputra28/belajar-gin-1/internal/users"

	"gorm.io/gorm"
)
type Contact struct {
	ID		uint		`gorm:"column:contact_id;primaryKey" json:"id"`
	UserID    uint       `gorm:"not null;index" json:"user_id"`
	User      users.User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	FirstName string     `gorm:"type:varchar(255);not null" json:"first_name"`
	LastName  string     `gorm:"type:varchar(255)" json:"last_name"`
	Email     string     `gorm:"type:varchar(255);not null" json:"email"`
	Phone     string     `gorm:"type:varchar(255)" json:"phone"`
	CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Contact) TableName() string {
	return "contacts"
}
type CreateContactRequest struct {
	FirstName string `json:"first_name" binding:"required,min=2,max=100"`
	LastName  string `json:"last_name" binding:"omitempty,max=100"`
	Email     string `json:"email" binding:"required,email"`
	Phone     string `json:"phone" binding:"omitempty,max=20"`
}

type UpdateContactRequest struct {
	FirstName string `json:"first_name" binding:"omitempty,min=2,max=100"`
	LastName  string `json:"last_name" binding:"omitempty,max=100"`
	Email     string `json:"email" binding:"omitempty,email"`
	Phone     string `json:"phone" binding:"omitempty,max=20"`
}

type ContactResponse struct {
	ID        uint   `json:"id"`
	UserID    uint   `json:"user_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

type GetContactsResponse struct {
	Data       []Contact `json:"data"`
	Page       int       `json:"page"`
	Limit      int       `json:"limit"`
	Total      int       `json:"total"`
	TotalPages int       `json:"total_pages"`
}