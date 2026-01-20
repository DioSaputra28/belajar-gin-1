package addresses

import (
	"time"

	"github.com/DioSaputra28/belajar-gin-1/internal/contacts"

	"gorm.io/gorm"
)

type Address struct {
	ID         uint             `gorm:"column:address_id;primaryKey" json:"id"`
	ContactID  uint             `gorm:"not null;index" json:"contact_id"`
	Contact    contacts.Contact `gorm:"foreignKey:ContactID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	Street     string           `gorm:"type:varchar(255)" json:"street"`
	City       string           `gorm:"type:varchar(255)" json:"city"`
	State      string           `gorm:"type:varchar(255)" json:"state"`
	PostalCode string           `gorm:"type:varchar(20)" json:"postal_code"`
	Country    string           `gorm:"type:varchar(100);not null" json:"country"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
	DeletedAt gorm.DeletedAt   `gorm:"index" json:"-"`
}

func (Address) TableName() string {
	return "addresses"
}

type CreateAddressRequest struct {
	ContactID  uint   `json:"contact_id" binding:"required"`
	Street     string `json:"street" binding:"omitempty,max=255"`
	City       string `json:"city" binding:"omitempty,max=100"`
	State      string `json:"state" binding:"omitempty,max=100"`
	PostalCode string `json:"postal_code" binding:"omitempty,max=20"`
	Country    string `json:"country" binding:"required,max=100"`
}

type UpdateAddressRequest struct {
	Street     string `json:"street" binding:"omitempty,max=255"`
	City       string `json:"city" binding:"omitempty,max=100"`
	State      string `json:"state" binding:"omitempty,max=100"`
	PostalCode string `json:"postal_code" binding:"omitempty,max=20"`
	Country    string `json:"country" binding:"omitempty,max=100"`
}

type AddressResponse struct {
	ID         uint   `json:"id"`
	ContactID  uint   `json:"contact_id"`
	Street     string `json:"street"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country"`
}

type GetAddressesResponse struct {
	Addresses []Address `json:"addresses"`
	Page      int               `json:"page"`
	Limit     int               `json:"limit"`
	Total     int               `json:"total"`
	TotalPage int               `json:"total_page"`
}