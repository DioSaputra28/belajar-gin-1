package test

import (
	"testing"

	"github.com/DioSaputra28/belajar-gin-1/internal/contacts"
	"github.com/DioSaputra28/belajar-gin-1/internal/users"
)

// ========== Contacts Repository Integration Tests ==========

func TestContactRepository_GetContacts_Integration(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestDBAfter(t, db)

	repo := contacts.NewContactRepository(db)

	// Create test users
	user1 := users.User{Name: "User 1", Email: "user1@example.com", Password: "hash1"}
	user2 := users.User{Name: "User 2", Email: "user2@example.com", Password: "hash2"}
	db.Create(&user1)
	db.Create(&user2)

	// Create contacts for different users
	db.Create(&contacts.Contact{UserID: user1.ID, FirstName: "John", LastName: "Doe", Email: "john@example.com", Phone: "123"})
	db.Create(&contacts.Contact{UserID: user1.ID, FirstName: "Jane", LastName: "Smith", Email: "jane@example.com", Phone: "456"})
	db.Create(&contacts.Contact{UserID: user2.ID, FirstName: "Bob", LastName: "Wilson", Email: "bob@example.com", Phone: "789"})

	// Test get contacts for user1 only
	result, err := repo.GetContacts(1, 10, int(user1.ID), "")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// BUG CHECK: Should only return user1's contacts (2), not all contacts (3)
	if len(result.Data) != 2 {
		t.Errorf("BUG FOUND: Expected 2 contacts for user1, got %d (user_id filter not working!)", len(result.Data))
	}

	if result.Total != 2 {
		t.Errorf("BUG FOUND: Expected total 2 for user1, got %d (user_id filter not working!)", result.Total)
	}

	// Verify no contacts from user2 are returned
	for _, contact := range result.Data {
		if contact.UserID != user1.ID {
			t.Errorf("BUG FOUND: Contact belongs to user %d, not user %d (authorization breach!)", contact.UserID, user1.ID)
		}
	}
}

func TestContactRepository_GetContacts_WithSearch_Integration(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestDBAfter(t, db)

	repo := contacts.NewContactRepository(db)

	// Create test user
	user1 := users.User{Name: "User 1", Email: "user1@example.com", Password: "hash1"}
	db.Create(&user1)

	// Create contacts
	db.Create(&contacts.Contact{UserID: user1.ID, FirstName: "John", Email: "john@example.com"})
	db.Create(&contacts.Contact{UserID: user1.ID, FirstName: "Jane", Email: "jane@example.com"})

	// Test search
	result, err := repo.GetContacts(1, 10, int(user1.ID), "john")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(result.Data) != 1 {
		t.Errorf("Expected 1 contact, got %d", len(result.Data))
	}

	if result.Data[0].FirstName != "John" {
		t.Errorf("Expected 'John', got '%s'", result.Data[0].FirstName)
	}
}

func TestContactRepository_CreateContact_Integration(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestDBAfter(t, db)

	repo := contacts.NewContactRepository(db)

	// Create test user
	user1 := users.User{Name: "User 1", Email: "user1@example.com", Password: "hash1"}
	db.Create(&user1)

	// Create contact
	newContact := contacts.Contact{
		UserID:    user1.ID,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		Phone:     "1234567890",
	}

	response, err := repo.CreateContact(newContact)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if response.FirstName != "John" {
		t.Errorf("Expected 'John', got '%s'", response.FirstName)
	}

	// Verify in database
	var contact contacts.Contact
	db.Where("email = ?", "john@example.com").First(&contact)
	if contact.ID == 0 {
		t.Error("Contact was not saved to database")
	}
}

func TestContactRepository_FindContactById_Integration(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestDBAfter(t, db)

	repo := contacts.NewContactRepository(db)

	// Create test user
	user1 := users.User{Name: "User 1", Email: "user1@example.com", Password: "hash1"}
	db.Create(&user1)

	// Create contact
	createdContact := contacts.Contact{
		UserID:    user1.ID,
		FirstName: "John",
		Email:     "john@example.com",
	}
	db.Create(&createdContact)

	// Find contact
	contact, err := repo.FindContactById(createdContact.ID, user1.ID)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if contact.ID != createdContact.ID {
		t.Errorf("Expected ID %d, got %d", createdContact.ID, contact.ID)
	}
}

func TestContactRepository_FindContactById_WrongUser_Integration(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestDBAfter(t, db)

	repo := contacts.NewContactRepository(db)

	// Create test users
	user1 := users.User{Name: "User 1", Email: "user1@example.com", Password: "hash1"}
	user2 := users.User{Name: "User 2", Email: "user2@example.com", Password: "hash2"}
	db.Create(&user1)
	db.Create(&user2)

	// Create contact for user1
	createdContact := contacts.Contact{
		UserID:    user1.ID,
		FirstName: "John",
		Email:     "john@example.com",
	}
	db.Create(&createdContact)

	// Try to find with user2's ID (should fail - authorization check)
	_, err := repo.FindContactById(createdContact.ID, user2.ID)

	if err == nil {
		t.Error("BUG FOUND: Should not find contact belonging to different user (authorization breach!)")
	}
}

func TestContactRepository_UpdateContact_Integration(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestDBAfter(t, db)

	repo := contacts.NewContactRepository(db)

	// Create test user
	user1 := users.User{Name: "User 1", Email: "user1@example.com", Password: "hash1"}
	db.Create(&user1)

	// Create contact
	createdContact := contacts.Contact{
		UserID:    user1.ID,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		Phone:     "123",
	}
	db.Create(&createdContact)

	// Update contact
	updateContact := &contacts.Contact{
		FirstName: "John Updated",
		Email:     "john.updated@example.com",
	}

	err := repo.UpdateContact(createdContact.ID, user1.ID, updateContact)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify update in database
	var contact contacts.Contact
	db.First(&contact, createdContact.ID)
	if contact.FirstName != "John Updated" {
		t.Errorf("Expected 'John Updated', got '%s'", contact.FirstName)
	}
}

func TestContactRepository_UpdateContact_PartialUpdate_Integration(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestDBAfter(t, db)

	repo := contacts.NewContactRepository(db)

	// Create test user
	user1 := users.User{Name: "User 1", Email: "user1@example.com", Password: "hash1"}
	db.Create(&user1)

	// Create contact
	createdContact := contacts.Contact{
		UserID:    user1.ID,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		Phone:     "123",
	}
	db.Create(&createdContact)

	// Update only first name
	updateContact := &contacts.Contact{
		FirstName: "John Updated",
	}

	err := repo.UpdateContact(createdContact.ID, user1.ID, updateContact)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify update in database
	var contact contacts.Contact
	db.First(&contact, createdContact.ID)

	// BUG CHECK: Other fields should NOT be cleared
	if contact.Email == "" {
		t.Error("BUG FOUND: Email was cleared when it should remain unchanged!")
	}
	if contact.LastName == "" {
		t.Error("BUG FOUND: LastName was cleared when it should remain unchanged!")
	}
	if contact.Phone == "" {
		t.Error("BUG FOUND: Phone was cleared when it should remain unchanged!")
	}
}

func TestContactRepository_DeleteContact_Integration(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestDBAfter(t, db)

	repo := contacts.NewContactRepository(db)

	// Create test user
	user1 := users.User{Name: "User 1", Email: "user1@example.com", Password: "hash1"}
	db.Create(&user1)

	// Create contact
	createdContact := contacts.Contact{
		UserID:    user1.ID,
		FirstName: "John",
		Email:     "john@example.com",
	}
	db.Create(&createdContact)

	// Delete contact
	err := repo.DeleteContact(createdContact.ID, user1.ID)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify deletion
	var contact contacts.Contact
	result := db.First(&contact, createdContact.ID)
	if result.Error == nil {
		t.Error("Contact should be deleted but was still found")
	}
}

func TestContactRepository_DeleteContact_WrongUser_Integration(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestDBAfter(t, db)

	repo := contacts.NewContactRepository(db)

	// Create test users
	user1 := users.User{Name: "User 1", Email: "user1@example.com", Password: "hash1"}
	user2 := users.User{Name: "User 2", Email: "user2@example.com", Password: "hash2"}
	db.Create(&user1)
	db.Create(&user2)

	// Create contact for user1
	createdContact := contacts.Contact{
		UserID:    user1.ID,
		FirstName: "John",
		Email:     "john@example.com",
	}
	db.Create(&createdContact)

	// Try to delete with user2's ID
	err := repo.DeleteContact(createdContact.ID, user2.ID)

	// Should succeed (no error) but contact should NOT be deleted
	if err != nil {
		t.Logf("Delete returned error: %v", err)
	}

	// BUG CHECK: Verify contact still exists (should not be deleted by wrong user)
	var contact contacts.Contact
	result := db.First(&contact, createdContact.ID)
	if result.Error != nil {
		t.Error("BUG FOUND: Contact was deleted by wrong user (authorization breach!)")
	}
}
