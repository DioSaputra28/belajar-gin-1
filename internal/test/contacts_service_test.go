package test

import (
	"errors"
	"testing"

	"github.com/DioSaputra28/belajar-gin-1/internal/contacts"
	"gorm.io/gorm"
)

// ========== GetContacts Tests ==========

// TestGetContacts_Success tests successful retrieval with pagination
func TestGetContacts_Success(t *testing.T) {
	mockRepo := &MockContactRepository{
		GetContactsFunc: func(page, limit, user_id int, search string) (*contacts.GetContactsResponse, error) {
			return &contacts.GetContactsResponse{
				Data: []contacts.Contact{
					{ID: 1, UserID: uint(user_id), FirstName: "John", LastName: "Doe", Email: "john@example.com", Phone: "1234567890"},
					{ID: 2, UserID: uint(user_id), FirstName: "Jane", LastName: "Smith", Email: "jane@example.com", Phone: "0987654321"},
				},
				Page:       page,
				Limit:      limit,
				Total:      2,
				TotalPages: 1,
			}, nil
		},
	}

	service := contacts.NewContactService(mockRepo)

	result, err := service.GetContacts(1, 10, 1, "")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(result.Data) != 2 {
		t.Errorf("Expected 2 contacts, got %d", len(result.Data))
	}

	if result.Data[0].FirstName != "John" {
		t.Errorf("Expected first contact name 'John', got '%s'", result.Data[0].FirstName)
	}

	if result.Total != 2 {
		t.Errorf("Expected total 2, got %d", result.Total)
	}
}

// TestGetContacts_WithSearch tests contact retrieval with search query
func TestGetContacts_WithSearch(t *testing.T) {
	mockRepo := &MockContactRepository{
		GetContactsFunc: func(page, limit, user_id int, search string) (*contacts.GetContactsResponse, error) {
			if search == "john" {
				return &contacts.GetContactsResponse{
					Data: []contacts.Contact{
						{ID: 1, UserID: uint(user_id), FirstName: "John", Email: "john@example.com"},
					},
					Page:       page,
					Limit:      limit,
					Total:      1,
					TotalPages: 1,
				}, nil
			}
			return &contacts.GetContactsResponse{
				Data:       []contacts.Contact{},
				Page:       page,
				Limit:      limit,
				Total:      0,
				TotalPages: 0,
			}, nil
		},
	}

	service := contacts.NewContactService(mockRepo)

	result, err := service.GetContacts(1, 10, 1, "john")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(result.Data) != 1 {
		t.Errorf("Expected 1 contact, got %d", len(result.Data))
	}

	if result.Data[0].Email != "john@example.com" {
		t.Errorf("Expected email 'john@example.com', got '%s'", result.Data[0].Email)
	}
}

// TestGetContacts_EmptyResult tests contact retrieval with no results
func TestGetContacts_EmptyResult(t *testing.T) {
	mockRepo := &MockContactRepository{
		GetContactsFunc: func(page, limit, user_id int, search string) (*contacts.GetContactsResponse, error) {
			return &contacts.GetContactsResponse{
				Data:       []contacts.Contact{},
				Page:       page,
				Limit:      limit,
				Total:      0,
				TotalPages: 0,
			}, nil
		},
	}

	service := contacts.NewContactService(mockRepo)

	result, err := service.GetContacts(1, 10, 1, "nonexistent")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(result.Data) != 0 {
		t.Errorf("Expected 0 contacts, got %d", len(result.Data))
	}
}

// TestGetContacts_DatabaseError tests contact retrieval with database error
func TestGetContacts_DatabaseError(t *testing.T) {
	mockRepo := &MockContactRepository{
		GetContactsFunc: func(page, limit, user_id int, search string) (*contacts.GetContactsResponse, error) {
			return nil, errors.New("database connection error")
		},
	}

	service := contacts.NewContactService(mockRepo)

	_, err := service.GetContacts(1, 10, 1, "")

	if err == nil {
		t.Error("Expected database error, got nil")
	}
}

// ========== CreateContact Tests ==========

// TestCreateContact_Success tests successful contact creation
func TestCreateContact_Success(t *testing.T) {
	mockRepo := &MockContactRepository{
		CreateContactFunc: func(contact contacts.Contact) (*contacts.ContactResponse, error) {
			return &contacts.ContactResponse{
				ID:        1,
				UserID:    contact.UserID,
				FirstName: contact.FirstName,
				LastName:  contact.LastName,
				Email:     contact.Email,
				Phone:     contact.Phone,
			}, nil
		},
	}

	service := contacts.NewContactService(mockRepo)

	newContact := contacts.Contact{
		UserID:    1,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		Phone:     "1234567890",
	}

	response, err := service.CreateContact(newContact)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if response.FirstName != "John" {
		t.Errorf("Expected first name 'John', got '%s'", response.FirstName)
	}

	if response.Email != "john@example.com" {
		t.Errorf("Expected email 'john@example.com', got '%s'", response.Email)
	}
}

// TestCreateContact_MinimumData tests contact creation with minimum data
func TestCreateContact_MinimumData(t *testing.T) {
	mockRepo := &MockContactRepository{
		CreateContactFunc: func(contact contacts.Contact) (*contacts.ContactResponse, error) {
			return &contacts.ContactResponse{
				ID:        1,
				UserID:    contact.UserID,
				FirstName: contact.FirstName,
				Email:     contact.Email,
			}, nil
		},
	}

	service := contacts.NewContactService(mockRepo)

	newContact := contacts.Contact{
		UserID:    1,
		FirstName: "Jo",
		Email:     "jo@example.com",
	}

	response, err := service.CreateContact(newContact)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if response.FirstName != "Jo" {
		t.Errorf("Expected first name 'Jo', got '%s'", response.FirstName)
	}
}

// TestCreateContact_DatabaseError tests contact creation with database error
func TestCreateContact_DatabaseError(t *testing.T) {
	mockRepo := &MockContactRepository{
		CreateContactFunc: func(contact contacts.Contact) (*contacts.ContactResponse, error) {
			return nil, errors.New("duplicate key value violates unique constraint")
		},
	}

	service := contacts.NewContactService(mockRepo)

	newContact := contacts.Contact{
		UserID:    1,
		FirstName: "John",
		Email:     "john@example.com",
	}

	_, err := service.CreateContact(newContact)

	if err == nil {
		t.Error("Expected database error, got nil")
	}
}

// ========== FindContactById Tests ==========

// TestFindContactById_Success tests successful contact retrieval by ID
func TestFindContactById_Success(t *testing.T) {
	mockRepo := &MockContactRepository{
		FindContactByIdFunc: func(id, user_id uint) (*contacts.Contact, error) {
			return &contacts.Contact{
				ID: id,
				UserID:    user_id,
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john@example.com",
				Phone:     "1234567890",
			}, nil
		},
	}

	service := contacts.NewContactService(mockRepo)

	contact, err := service.FindContactById(1, 1)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if contact.ID != 1 {
		t.Errorf("Expected contact ID 1, got %d", contact.ID)
	}

	if contact.FirstName != "John" {
		t.Errorf("Expected first name 'John', got '%s'", contact.FirstName)
	}
}

// TestFindContactById_NotFound tests finding non-existent contact
func TestFindContactById_NotFound(t *testing.T) {
	mockRepo := &MockContactRepository{
		FindContactByIdFunc: func(id, user_id uint) (*contacts.Contact, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}

	service := contacts.NewContactService(mockRepo)

	_, err := service.FindContactById(999, 1)

	if err == nil {
		t.Error("Expected error for non-existent contact, got nil")
	}

	if err.Error() != "contact not found" {
		t.Errorf("Expected 'contact not found' error, got '%s'", err.Error())
	}
}

// TestFindContactById_WrongUser tests finding contact with wrong user_id (authorization)
func TestFindContactById_WrongUser(t *testing.T) {
	mockRepo := &MockContactRepository{
		FindContactByIdFunc: func(id, user_id uint) (*contacts.Contact, error) {
			// Simulate contact not found because user_id doesn't match
			return nil, gorm.ErrRecordNotFound
		},
	}

	service := contacts.NewContactService(mockRepo)

	_, err := service.FindContactById(1, 999)

	if err == nil {
		t.Error("Expected error for unauthorized access, got nil")
	}

	if err.Error() != "contact not found" {
		t.Errorf("Expected 'contact not found' error, got '%s'", err.Error())
	}
}

// TestFindContactById_DatabaseError tests finding contact with database error
func TestFindContactById_DatabaseError(t *testing.T) {
	mockRepo := &MockContactRepository{
		FindContactByIdFunc: func(id, user_id uint) (*contacts.Contact, error) {
			return nil, errors.New("database connection error")
		},
	}

	service := contacts.NewContactService(mockRepo)

	_, err := service.FindContactById(1, 1)

	if err == nil {
		t.Error("Expected database error, got nil")
	}

	// Should propagate the database error, not convert to "contact not found"
	if err.Error() == "contact not found" {
		t.Error("Database error should be propagated, not converted to 'contact not found'")
	}
}

// ========== UpdateContact Tests ==========

// TestUpdateContact_Success tests successful contact update
func TestUpdateContact_Success(t *testing.T) {
	mockRepo := &MockContactRepository{
		FindContactByIdFunc: func(id, user_id uint) (*contacts.Contact, error) {
			return &contacts.Contact{
				ID: id,
				UserID:    user_id,
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john@example.com",
				Phone:     "1234567890",
			}, nil
		},
		UpdateContactFunc: func(id, user_id uint, contact *contacts.Contact) error {
			return nil
		},
	}

	service := contacts.NewContactService(mockRepo)

	updateContact := contacts.Contact{
		FirstName: "John Updated",
		Email:     "john.updated@example.com",
	}

	err := service.UpdateContact(1, 1, updateContact)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

// TestUpdateContact_FirstNameOnly tests updating only first name
func TestUpdateContact_FirstNameOnly(t *testing.T) {
	mockRepo := &MockContactRepository{
		FindContactByIdFunc: func(id, user_id uint) (*contacts.Contact, error) {
			return &contacts.Contact{
				ID: id,
				UserID:    user_id,
				FirstName: "John",
				Email:     "john@example.com",
			}, nil
		},
		UpdateContactFunc: func(id, user_id uint, contact *contacts.Contact) error {
			return nil
		},
	}

	service := contacts.NewContactService(mockRepo)

	updateContact := contacts.Contact{
		FirstName: "John Updated",
	}

	err := service.UpdateContact(1, 1, updateContact)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

// TestUpdateContact_EmailOnly tests updating only email
func TestUpdateContact_EmailOnly(t *testing.T) {
	mockRepo := &MockContactRepository{
		FindContactByIdFunc: func(id, user_id uint) (*contacts.Contact, error) {
			return &contacts.Contact{
				ID: id,
				UserID:    user_id,
				FirstName: "John",
				Email:     "john@example.com",
			}, nil
		},
		UpdateContactFunc: func(id, user_id uint, contact *contacts.Contact) error {
			return nil
		},
	}

	service := contacts.NewContactService(mockRepo)

	updateContact := contacts.Contact{
		Email: "john.new@example.com",
	}

	err := service.UpdateContact(1, 1, updateContact)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

// TestUpdateContact_NotFound tests updating non-existent contact
func TestUpdateContact_NotFound(t *testing.T) {
	mockRepo := &MockContactRepository{
		FindContactByIdFunc: func(id, user_id uint) (*contacts.Contact, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}

	service := contacts.NewContactService(mockRepo)

	updateContact := contacts.Contact{
		FirstName: "John Updated",
	}

	err := service.UpdateContact(999, 1, updateContact)

	if err == nil {
		t.Error("Expected error for non-existent contact, got nil")
	}

	if err.Error() != "contact not found" {
		t.Errorf("Expected 'contact not found' error, got '%s'", err.Error())
	}
}

// TestUpdateContact_WrongUser tests updating contact with wrong user_id (authorization)
func TestUpdateContact_WrongUser(t *testing.T) {
	mockRepo := &MockContactRepository{
		FindContactByIdFunc: func(id, user_id uint) (*contacts.Contact, error) {
			// Simulate contact not found because user_id doesn't match
			return nil, gorm.ErrRecordNotFound
		},
	}

	service := contacts.NewContactService(mockRepo)

	updateContact := contacts.Contact{
		FirstName: "John Updated",
	}

	err := service.UpdateContact(1, 999, updateContact)

	if err == nil {
		t.Error("Expected error for unauthorized access, got nil")
	}

	if err.Error() != "contact not found" {
		t.Errorf("Expected 'contact not found' error, got '%s'", err.Error())
	}
}

// TestUpdateContact_DatabaseError tests update with database error
func TestUpdateContact_DatabaseError(t *testing.T) {
	mockRepo := &MockContactRepository{
		FindContactByIdFunc: func(id, user_id uint) (*contacts.Contact, error) {
			return &contacts.Contact{
				ID: id,
				UserID:    user_id,
				FirstName: "John",
				Email:     "john@example.com",
			}, nil
		},
		UpdateContactFunc: func(id, user_id uint, contact *contacts.Contact) error {
			return errors.New("database connection error")
		},
	}

	service := contacts.NewContactService(mockRepo)

	updateContact := contacts.Contact{
		FirstName: "John Updated",
	}

	err := service.UpdateContact(1, 1, updateContact)

	if err == nil {
		t.Error("Expected database error, got nil")
	}
}

// ========== DeleteContact Tests ==========

// TestDeleteContact_Success tests successful contact deletion
func TestDeleteContact_Success(t *testing.T) {
	mockRepo := &MockContactRepository{
		FindContactByIdFunc: func(id, user_id uint) (*contacts.Contact, error) {
			return &contacts.Contact{
				ID: id,
				UserID:    user_id,
				FirstName: "John",
				Email:     "john@example.com",
			}, nil
		},
		DeleteContactFunc: func(id, user_id uint) error {
			return nil
		},
	}

	service := contacts.NewContactService(mockRepo)

	err := service.DeleteContact(1, 1)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

// TestDeleteContact_NotFound tests deleting non-existent contact
func TestDeleteContact_NotFound(t *testing.T) {
	mockRepo := &MockContactRepository{
		FindContactByIdFunc: func(id, user_id uint) (*contacts.Contact, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}

	service := contacts.NewContactService(mockRepo)

	err := service.DeleteContact(999, 1)

	if err == nil {
		t.Error("Expected error for non-existent contact, got nil")
	}

	if err.Error() != "contact not found" {
		t.Errorf("Expected 'contact not found' error, got '%s'", err.Error())
	}
}

// TestDeleteContact_WrongUser tests deleting contact with wrong user_id (authorization)
func TestDeleteContact_WrongUser(t *testing.T) {
	mockRepo := &MockContactRepository{
		FindContactByIdFunc: func(id, user_id uint) (*contacts.Contact, error) {
			// Simulate contact not found because user_id doesn't match
			return nil, gorm.ErrRecordNotFound
		},
	}

	service := contacts.NewContactService(mockRepo)

	err := service.DeleteContact(1, 999)

	if err == nil {
		t.Error("Expected error for unauthorized access, got nil")
	}

	if err.Error() != "contact not found" {
		t.Errorf("Expected 'contact not found' error, got '%s'", err.Error())
	}
}

// TestDeleteContact_DatabaseError tests deletion with database error
func TestDeleteContact_DatabaseError(t *testing.T) {
	mockRepo := &MockContactRepository{
		FindContactByIdFunc: func(id, user_id uint) (*contacts.Contact, error) {
			return &contacts.Contact{
				ID: id,
				UserID:    user_id,
				FirstName: "John",
				Email:     "john@example.com",
			}, nil
		},
		DeleteContactFunc: func(id, user_id uint) error {
			return errors.New("database connection error")
		},
	}

	service := contacts.NewContactService(mockRepo)

	err := service.DeleteContact(1, 1)

	if err == nil {
		t.Error("Expected database error, got nil")
	}
}
