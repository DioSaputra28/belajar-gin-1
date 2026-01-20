package test

import (
	"github.com/DioSaputra28/belajar-gin-1/internal/addresses"
	"github.com/DioSaputra28/belajar-gin-1/internal/contacts"
	"github.com/DioSaputra28/belajar-gin-1/internal/users"
	"gorm.io/gorm"
)

// MockUserRepository is a mock implementation of users.UserRepository
type MockUserRepository struct {
	GetUsersFunc     func(page, limit int, search string) (*users.GetUsersResponse, error)
	CreateUserFunc   func(user users.CreateUserRequest) (*users.CreateUserResonse, error)
	UpdateUserFunc   func(id uint, user users.UpdateUserRequest) error
	FindUserByIdFunc func(id uint) (*users.UserResponse, error)
	DeleteUserFunc   func(id uint) error
}

// GetUsers implements users.UserRepository
func (m *MockUserRepository) GetUsers(page, limit int, search string) (*users.GetUsersResponse, error) {
	if m.GetUsersFunc != nil {
		return m.GetUsersFunc(page, limit, search)
	}
	return nil, nil
}

// CreateUser implements users.UserRepository
func (m *MockUserRepository) CreateUser(user users.CreateUserRequest) (*users.CreateUserResonse, error) {
	if m.CreateUserFunc != nil {
		return m.CreateUserFunc(user)
	}
	return nil, nil
}

// UpdateUser implements users.UserRepository
func (m *MockUserRepository) UpdateUser(id uint, user users.UpdateUserRequest) error {
	if m.UpdateUserFunc != nil {
		return m.UpdateUserFunc(id, user)
	}
	return nil
}

// FindUserById implements users.UserRepository
func (m *MockUserRepository) FindUserById(id uint) (*users.UserResponse, error) {
	if m.FindUserByIdFunc != nil {
		return m.FindUserByIdFunc(id)
	}
	return nil, nil
}

// DeleteUser implements users.UserRepository
func (m *MockUserRepository) DeleteUser(id uint) error {
	if m.DeleteUserFunc != nil {
		return m.DeleteUserFunc(id)
	}
	return nil
}

// Helper function to create a sample user
func CreateSampleUser(id uint, name, email, password string) *users.User {
	return &users.User{
		ID:       id,
		Name:     name,
		Email:    email,
		Password: password,
	}
}

// Helper function to check if error is ErrRecordNotFound
func IsRecordNotFoundError(err error) bool {
	return err == gorm.ErrRecordNotFound
}

// MockContactRepository is a mock implementation of contacts.ContactRepository
type MockContactRepository struct {
	GetContactsFunc     func(page, limit, user_id int, search string) (*contacts.GetContactsResponse, error)
	CreateContactFunc   func(contact contacts.Contact) (*contacts.ContactResponse, error)
	FindContactByIdFunc func(id, user_id uint) (*contacts.Contact, error)
	UpdateContactFunc   func(id, user_id uint, contact *contacts.Contact) error
	DeleteContactFunc   func(id, user_id uint) error
}

// GetContacts implements contacts.ContactRepository
func (m *MockContactRepository) GetContacts(page, limit, user_id int, search string) (*contacts.GetContactsResponse, error) {
	if m.GetContactsFunc != nil {
		return m.GetContactsFunc(page, limit, user_id, search)
	}
	return nil, nil
}

// CreateContact implements contacts.ContactRepository
func (m *MockContactRepository) CreateContact(contact contacts.Contact) (*contacts.ContactResponse, error) {
	if m.CreateContactFunc != nil {
		return m.CreateContactFunc(contact)
	}
	return nil, nil
}

// FindContactById implements contacts.ContactRepository
func (m *MockContactRepository) FindContactById(id, user_id uint) (*contacts.Contact, error) {
	if m.FindContactByIdFunc != nil {
		return m.FindContactByIdFunc(id, user_id)
	}
	return nil, nil
}

// UpdateContact implements contacts.ContactRepository
func (m *MockContactRepository) UpdateContact(id, user_id uint, contact *contacts.Contact) error {
	if m.UpdateContactFunc != nil {
		return m.UpdateContactFunc(id, user_id, contact)
	}
	return nil
}

// DeleteContact implements contacts.ContactRepository
func (m *MockContactRepository) DeleteContact(id, user_id uint) error {
	if m.DeleteContactFunc != nil {
		return m.DeleteContactFunc(id, user_id)
	}
	return nil
}

// MockAddressRepository is a mock implementation of addresses.AddressRepository
type MockAddressRepository struct {
	CreateAddressFunc   func(address addresses.CreateAddressRequest) (*addresses.AddressResponse, error)
	GetAddressesFunc    func(contact_id uint, page int, limit int, search string) (*addresses.GetAddressesResponse, error)
	UpdateAddressFunc   func(address_id uint, address *addresses.Address) (*addresses.AddressResponse, error)
	FindAddressByIdFunc func(address_id uint) (*addresses.Address, error)
	DeleteAddressFunc   func(address_id uint) error
	FindContactByIdFunc func(id, user_id uint) (*contacts.Contact, error)
}

// CreateAddress implements addresses.AddressRepository
func (m *MockAddressRepository) CreateAddress(address addresses.CreateAddressRequest) (*addresses.AddressResponse, error) {
	if m.CreateAddressFunc != nil {
		return m.CreateAddressFunc(address)
	}
	return nil, nil
}

// GetAddresses implements addresses.AddressRepository
func (m *MockAddressRepository) GetAddresses(contact_id uint, page int, limit int, search string) (*addresses.GetAddressesResponse, error) {
	if m.GetAddressesFunc != nil {
		return m.GetAddressesFunc(contact_id, page, limit, search)
	}
	return nil, nil
}

// UpdateAddress implements addresses.AddressRepository
func (m *MockAddressRepository) UpdateAddress(address_id uint, address *addresses.Address) (*addresses.AddressResponse, error) {
	if m.UpdateAddressFunc != nil {
		return m.UpdateAddressFunc(address_id, address)
	}
	return nil, nil
}

// FindAddressById implements addresses.AddressRepository
func (m *MockAddressRepository) FindAddressById(address_id uint) (*addresses.Address, error) {
	if m.FindAddressByIdFunc != nil {
		return m.FindAddressByIdFunc(address_id)
	}
	return nil, nil
}

// DeleteAddress implements addresses.AddressRepository
func (m *MockAddressRepository) DeleteAddress(address_id uint) error {
	if m.DeleteAddressFunc != nil {
		return m.DeleteAddressFunc(address_id)
	}
	return nil
}

// FindContactById implements addresses.AddressRepository
func (m *MockAddressRepository) FindContactById(id, user_id uint) (*contacts.Contact, error) {
	if m.FindContactByIdFunc != nil {
		return m.FindContactByIdFunc(id, user_id)
	}
	return nil, nil
}
