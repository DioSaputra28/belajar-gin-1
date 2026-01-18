package addresses

import (
	"github.com/DioSaputra28/belajar-gin-1/internal/contacts"

	"gorm.io/gorm"
)

type Address struct {
	gorm.Model
	ContactID  uint             `gorm:"not null;index" json:"contact_id"`
	Contact    contacts.Contact `gorm:"foreignKey:ContactID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	Street     string           `gorm:"type:varchar(255)" json:"street"`
	City       string           `gorm:"type:varchar(255)" json:"city"`
	State      string           `gorm:"type:varchar(255)" json:"state"`
	PostalCode string           `gorm:"type:varchar(20)" json:"postal_code"`
	Country    string           `gorm:"type:varchar(100);not null" json:"country"`
}

func (Address) TableName() string {
	return "addresses"
}

type CreateAddressRequest struct {
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
