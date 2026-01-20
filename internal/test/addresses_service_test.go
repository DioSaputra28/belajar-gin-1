package test

import (
	"errors"
	"testing"

	"github.com/DioSaputra28/belajar-gin-1/internal/addresses"
	"github.com/DioSaputra28/belajar-gin-1/internal/contacts"
	"gorm.io/gorm"
)

// ========== CreateAddress Tests ==========

func TestCreateAddress_Success(t *testing.T) {
	mockRepo := &MockAddressRepository{
		FindContactByIdFunc: func(id, user_id uint) (*contacts.Contact, error) {
			return &contacts.Contact{
				ID:        id,
				UserID:    user_id,
				FirstName: "John",
				Email:     "john@example.com",
			}, nil
		},
		CreateAddressFunc: func(address addresses.CreateAddressRequest) (*addresses.AddressResponse, error) {
			return &addresses.AddressResponse{
				ID:         1,
				ContactID:  address.ContactID,
				Street:     address.Street,
				City:       address.City,
				State:      address.State,
				PostalCode: address.PostalCode,
				Country:    address.Country,
			}, nil
		},
	}

	service := addresses.NewAddressService(mockRepo)

	request := addresses.CreateAddressRequest{
		ContactID:  1,
		Street:     "123 Main St",
		City:       "New York",
		State:      "NY",
		PostalCode: "10001",
		Country:    "USA",
	}

	response, err := service.CreateAddress(1, request)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if response.Street != "123 Main St" {
		t.Errorf("Expected street '123 Main St', got '%s'", response.Street)
	}
}

func TestCreateAddress_ContactNotFound(t *testing.T) {
	mockRepo := &MockAddressRepository{
		FindContactByIdFunc: func(id, user_id uint) (*contacts.Contact, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}

	service := addresses.NewAddressService(mockRepo)

	request := addresses.CreateAddressRequest{
		ContactID: 999,
		Country:   "USA",
	}

	_, err := service.CreateAddress(1, request)

	if err == nil {
		t.Error("Expected error for non-existent contact, got nil")
	}

	if err.Error() != "contact not found" {
		t.Errorf("Expected 'contact not found' error, got '%s'", err.Error())
	}
}

func TestCreateAddress_DatabaseError(t *testing.T) {
	mockRepo := &MockAddressRepository{
		FindContactByIdFunc: func(id, user_id uint) (*contacts.Contact, error) {
			return &contacts.Contact{ID: id, UserID: user_id}, nil
		},
		CreateAddressFunc: func(address addresses.CreateAddressRequest) (*addresses.AddressResponse, error) {
			return nil, errors.New("database connection error")
		},
	}

	service := addresses.NewAddressService(mockRepo)

	request := addresses.CreateAddressRequest{
		ContactID: 1,
		Country:   "USA",
	}

	_, err := service.CreateAddress(1, request)

	if err == nil {
		t.Error("Expected database error, got nil")
	}
}

// ========== GetAddresses Tests ==========

func TestGetAddresses_Success(t *testing.T) {
	mockRepo := &MockAddressRepository{
		GetAddressesFunc: func(contact_id uint, page int, limit int, search string) (*addresses.GetAddressesResponse, error) {
			return &addresses.GetAddressesResponse{
				Addresses: []addresses.Address{},
				Page:      page,
				Limit:     limit,
				Total:     2,
				TotalPage: 1,
			}, nil
		},
	}

	service := addresses.NewAddressService(mockRepo)

	result, err := service.GetAddresses(1, 1, 1, 10, "")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(result.Addresses) != 2 {
		t.Errorf("Expected 2 addresses, got %d", len(result.Addresses))
	}

	if result.Total != 2 {
		t.Errorf("Expected total 2, got %d", result.Total)
	}
}

func TestGetAddresses_WithSearch(t *testing.T) {
	mockRepo := &MockAddressRepository{
		GetAddressesFunc: func(contact_id uint, page int, limit int, search string) (*addresses.GetAddressesResponse, error) {
			if search == "New York" {
				return &addresses.GetAddressesResponse{
					Addresses: []addresses.Address{},
					Page:      page,
					Limit:     limit,
					Total:     1,
					TotalPage: 1,
				}, nil
			}
			return &addresses.GetAddressesResponse{
				Addresses: []addresses.Address{},
				Page:      page,
				Limit:     limit,
				Total:     0,
				TotalPage: 0,
			}, nil
		},
	}

	service := addresses.NewAddressService(mockRepo)

	result, err := service.GetAddresses(1, 1, 1, 10, "New York")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(result.Addresses) != 1 {
		t.Errorf("Expected 1 address, got %d", len(result.Addresses))
	}
}

func TestGetAddresses_EmptyResult(t *testing.T) {
	mockRepo := &MockAddressRepository{
		GetAddressesFunc: func(contact_id uint, page int, limit int, search string) (*addresses.GetAddressesResponse, error) {
			return &addresses.GetAddressesResponse{
				Addresses: []addresses.Address{},
				Page:      page,
				Limit:     limit,
				Total:     0,
				TotalPage: 0,
			}, nil
		},
	}

	service := addresses.NewAddressService(mockRepo)

	result, err := service.GetAddresses(1, 1, 1, 10, "")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(result.Addresses) != 0 {
		t.Errorf("Expected 0 addresses, got %d", len(result.Addresses))
	}
}

func TestGetAddresses_DatabaseError(t *testing.T) {
	mockRepo := &MockAddressRepository{
		GetAddressesFunc: func(contact_id uint, page int, limit int, search string) (*addresses.GetAddressesResponse, error) {
			return nil, errors.New("database connection error")
		},
	}

	service := addresses.NewAddressService(mockRepo)

	_, err := service.GetAddresses(1, 1, 1, 10, "")

	if err == nil {
		t.Error("Expected database error, got nil")
	}
}

// ========== FindAddressById Tests ==========

func TestFindAddressById_Success(t *testing.T) {
	mockRepo := &MockAddressRepository{
		FindAddressByIdFunc: func(address_id uint) (*addresses.Address, error) {
			return &addresses.Address{
				ContactID:  1,
				Street:     "123 Main St",
				City:       "New York",
				State:      "NY",
				PostalCode: "10001",
				Country:    "USA",
			}, nil
		},
	}

	service := addresses.NewAddressService(mockRepo)

	address, err := service.FindAddressById(1, 1)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if address.ID != 1 {
		t.Errorf("Expected address ID 1, got %d", address.ID)
	}

	if address.City != "New York" {
		t.Errorf("Expected city 'New York', got '%s'", address.City)
	}
}

func TestFindAddressById_NotFound(t *testing.T) {
	mockRepo := &MockAddressRepository{
		FindAddressByIdFunc: func(address_id uint) (*addresses.Address, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}

	service := addresses.NewAddressService(mockRepo)

	_, err := service.FindAddressById(1, 999)

	if err == nil {
		t.Error("Expected error for non-existent address, got nil")
	}

	if err.Error() != "contact not found" {
		t.Errorf("Expected 'contact not found' error, got '%s'", err.Error())
	}
}

func TestFindAddressById_WrongUser(t *testing.T) {
	mockRepo := &MockAddressRepository{
		FindAddressByIdFunc: func(address_id uint) (*addresses.Address, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}

	service := addresses.NewAddressService(mockRepo)

	_, err := service.FindAddressById(999, 1)

	if err == nil {
		t.Error("Expected error for unauthorized access, got nil")
	}
}

func TestFindAddressById_DatabaseError(t *testing.T) {
	mockRepo := &MockAddressRepository{
		FindAddressByIdFunc: func(address_id uint) (*addresses.Address, error) {
			return nil, errors.New("database connection error")
		},
	}

	service := addresses.NewAddressService(mockRepo)

	_, err := service.FindAddressById(1, 1)

	if err == nil {
		t.Error("Expected database error, got nil")
	}

	if err.Error() == "contact not found" {
		t.Error("Database error should be propagated, not converted to 'contact not found'")
	}
}

// ========== UpdateAddress Tests ==========

func TestUpdateAddress_Success(t *testing.T) {
	mockRepo := &MockAddressRepository{
		FindAddressByIdFunc: func(address_id uint) (*addresses.Address, error) {
			return &addresses.Address{
				ContactID:  1,
				Street:     "123 Main St",
				City:       "New York",
				State:      "NY",
				PostalCode: "10001",
				Country:    "USA",
			}, nil
		},
		UpdateAddressFunc: func(address_id uint, address *addresses.Address) (*addresses.AddressResponse, error) {
			return &addresses.AddressResponse{
				ID:         address_id,
				ContactID:  address.ContactID,
				Street:     address.Street,
				City:       address.City,
				State:      address.State,
				PostalCode: address.PostalCode,
				Country:    address.Country,
			}, nil
		},
	}

	service := addresses.NewAddressService(mockRepo)

	request := addresses.UpdateAddressRequest{
		Street: "456 Updated St",
		City:   "Boston",
	}

	response, err := service.UpdateAddress(1, 1, request)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if response.Street != "456 Updated St" {
		t.Errorf("Expected street '456 Updated St', got '%s'", response.Street)
	}
}

func TestUpdateAddress_PartialUpdate(t *testing.T) {
	mockRepo := &MockAddressRepository{
		FindAddressByIdFunc: func(address_id uint) (*addresses.Address, error) {
			return &addresses.Address{
				ContactID: 1,
				Street:    "123 Main St",
				City:      "New York",
				Country:   "USA",
			}, nil
		},
		UpdateAddressFunc: func(address_id uint, address *addresses.Address) (*addresses.AddressResponse, error) {
			return &addresses.AddressResponse{
				ID:      address_id,
				Street:  address.Street,
				City:    address.City,
				Country: address.Country,
			}, nil
		},
	}

	service := addresses.NewAddressService(mockRepo)

	request := addresses.UpdateAddressRequest{
		City: "Boston",
	}

	response, err := service.UpdateAddress(1, 1, request)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if response.City != "Boston" {
		t.Errorf("Expected city 'Boston', got '%s'", response.City)
	}
}

func TestUpdateAddress_NotFound(t *testing.T) {
	mockRepo := &MockAddressRepository{
		FindAddressByIdFunc: func(address_id uint) (*addresses.Address, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}

	service := addresses.NewAddressService(mockRepo)

	request := addresses.UpdateAddressRequest{
		City: "Boston",
	}

	_, err := service.UpdateAddress(1, 999, request)

	if err == nil {
		t.Error("Expected error for non-existent address, got nil")
	}
}

func TestUpdateAddress_WrongUser(t *testing.T) {
	mockRepo := &MockAddressRepository{
		FindAddressByIdFunc: func(address_id uint) (*addresses.Address, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}

	service := addresses.NewAddressService(mockRepo)

	request := addresses.UpdateAddressRequest{
		City: "Boston",
	}

	_, err := service.UpdateAddress(999, 1, request)

	if err == nil {
		t.Error("Expected error for unauthorized access, got nil")
	}
}

func TestUpdateAddress_DatabaseError(t *testing.T) {
	mockRepo := &MockAddressRepository{
		FindAddressByIdFunc: func(address_id uint) (*addresses.Address, error) {
			return &addresses.Address{
				ContactID: 1,
				City:      "New York",
			}, nil
		},
		UpdateAddressFunc: func(address_id uint, address *addresses.Address) (*addresses.AddressResponse, error) {
			return nil, errors.New("database connection error")
		},
	}

	service := addresses.NewAddressService(mockRepo)

	request := addresses.UpdateAddressRequest{
		City: "Boston",
	}

	_, err := service.UpdateAddress(1, 1, request)

	if err == nil {
		t.Error("Expected database error, got nil")
	}
}

// ========== DeleteAddress Tests ==========

func TestDeleteAddress_Success(t *testing.T) {
	mockRepo := &MockAddressRepository{
		FindAddressByIdFunc: func(address_id uint) (*addresses.Address, error) {
			return &addresses.Address{
				ContactID: 1,
				City:      "New York",
			}, nil
		},
		DeleteAddressFunc: func(address_id uint) error {
			return nil
		},
	}

	service := addresses.NewAddressService(mockRepo)

	err := service.DeleteAddress(1, 1)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestDeleteAddress_NotFound(t *testing.T) {
	mockRepo := &MockAddressRepository{
		FindAddressByIdFunc: func(address_id uint) (*addresses.Address, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}

	service := addresses.NewAddressService(mockRepo)

	err := service.DeleteAddress(1, 999)

	if err == nil {
		t.Error("Expected error for non-existent address, got nil")
	}
}

func TestDeleteAddress_WrongUser(t *testing.T) {
	mockRepo := &MockAddressRepository{
		FindAddressByIdFunc: func(address_id uint) (*addresses.Address, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}

	service := addresses.NewAddressService(mockRepo)

	err := service.DeleteAddress(999, 1)

	if err == nil {
		t.Error("Expected error for unauthorized access, got nil")
	}
}

func TestDeleteAddress_DatabaseError(t *testing.T) {
	mockRepo := &MockAddressRepository{
		FindAddressByIdFunc: func(address_id uint) (*addresses.Address, error) {
			return &addresses.Address{
				ContactID: 1,
			}, nil
		},
		DeleteAddressFunc: func(address_id uint) error {
			return errors.New("database connection error")
		},
	}

	service := addresses.NewAddressService(mockRepo)

	err := service.DeleteAddress(1, 1)

	if err == nil {
		t.Error("Expected database error, got nil")
	}
}
