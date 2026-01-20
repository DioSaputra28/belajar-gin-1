package addresses

import (
	"github.com/DioSaputra28/belajar-gin-1/internal/contacts"

	"gorm.io/gorm"
)

type AddressRepository interface {
	CreateAddress(address CreateAddressRequest) (*AddressResponse, error)
	GetAddresses(contact_id uint, page int, limit int, search string) (*GetAddressesResponse, error)
	UpdateAddress(address_id uint, address *Address) (*AddressResponse, error)
	FindAddressById(address_id uint) (*Address, error)
	DeleteAddress(address_id uint) error
	FindContactById(id, user_id uint) (*contacts.Contact, error)
}

type addressRepository struct {
	db *gorm.DB
}

func NewAddressRepository(db *gorm.DB) AddressRepository {
	return &addressRepository{db: db}
}

func (a *addressRepository) CreateAddress(address CreateAddressRequest) (*AddressResponse, error) {
	address_db := Address{
		ContactID:  address.ContactID,
		Street:     address.Street,
		City:       address.City,
		State:      address.State,
		PostalCode: address.PostalCode,
		Country:    address.Country,
	}

	if err := a.db.Create(&address_db).Error; err != nil {
		return nil, err
	}
	return &AddressResponse{
		ID:         address_db.ID,
		ContactID:  address_db.ContactID,
		Street:     address_db.Street,
		City:       address_db.City,
		State:      address_db.State,
		PostalCode: address_db.PostalCode,
		Country:    address_db.Country,
	}, nil
}

func (a *addressRepository) GetAddresses(contact_id uint, page int, limit int, search string) (*GetAddressesResponse, error) {
	var addresses []Address
	var total int64

	// Build query dengan search filter
	query := a.db.Model(&Address{}).Where("contact_id = ?", contact_id)
	if search != "" {
		query = query.Where("street LIKE ? OR city LIKE ? OR state LIKE ? OR postal_code LIKE ? OR country LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// Count total records (sebelum pagination)
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// Get paginated data
	if err := query.Offset((page - 1) * limit).Limit(limit).Find(&addresses).Error; err != nil {
		return nil, err
	}

	// Hitung total pages
	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++ // Tambah 1 jika ada sisa
	}

	return &GetAddressesResponse{
		Addresses: addresses,
		Page:      page,
		Limit:     limit,
		Total:     int(total),
		TotalPage: totalPages,
	}, nil
}

func (a *addressRepository) UpdateAddress(address_id uint, address *Address) (*AddressResponse, error) {
	// Set ID to ensure update, not create
	address.ID = address_id

	if err := a.db.Save(&address).Error; err != nil {
		return nil, err
	}
	return &AddressResponse{
		ID:         address.ID,
		ContactID:  address.ContactID,
		Street:     address.Street,
		City:       address.City,
		State:      address.State,
		PostalCode: address.PostalCode,
		Country:    address.Country,
	}, nil
}

func (a *addressRepository) FindAddressById(address_id uint) (*Address, error) {
	var address Address
	if err := a.db.Where("address_id = ?", address_id).First(&address).Error; err != nil {
		return nil, err
	}
	return &address, nil
}

func (a *addressRepository) DeleteAddress(address_id uint) error {
	if err := a.db.Where("address_id = ?", address_id).Delete(&Address{}).Error; err != nil {
		return err
	}
	return nil
}

func (a *addressRepository) FindContactById(id uint, user_id uint) (*contacts.Contact, error) {
	contact_db := contacts.Contact{}
	err := a.db.Where("contact_id = ? AND user_id = ?", id, user_id).First(&contact_db).Error
	if err != nil {
		return nil, err
	}
	return &contact_db, nil
}
