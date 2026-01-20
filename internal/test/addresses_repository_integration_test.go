package test

import (
	"testing"

	"github.com/DioSaputra28/belajar-gin-1/internal/addresses"
	"github.com/DioSaputra28/belajar-gin-1/internal/contacts"
	"github.com/DioSaputra28/belajar-gin-1/internal/users"
)

// ========== Addresses Repository Integration Tests ==========

func TestAddressRepository_CreateAddress_Integration(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestDBAfter(t, db)

	repo := addresses.NewAddressRepository(db)

	// Create test user and contact
	user1 := users.User{Name: "User 1", Email: "user1@example.com", Password: "hash1"}
	db.Create(&user1)

	contact1 := contacts.Contact{UserID: user1.ID, FirstName: "John", Email: "john@example.com"}
	db.Create(&contact1)

	// Create address
	request := addresses.CreateAddressRequest{
		ContactID:  contact1.ID,
		Street:     "123 Main St",
		City:       "New York",
		State:      "NY",
		PostalCode: "10001",
		Country:    "USA",
	}

	response, err := repo.CreateAddress(request)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if response.Street != "123 Main St" {
		t.Errorf("Expected street '123 Main St', got '%s'", response.Street)
	}

	// Verify in database
	var address addresses.Address
	db.Where("contact_id = ?", contact1.ID).First(&address)
	if address.ID == 0 {
		t.Error("Address was not saved to database")
	}
}

func TestAddressRepository_CreateAddress_ForeignKeyConstraint_Integration(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestDBAfter(t, db)

	repo := addresses.NewAddressRepository(db)

	// Try to create address with non-existent contact
	request := addresses.CreateAddressRequest{
		ContactID: 999,
		Country:   "USA",
	}

	_, err := repo.CreateAddress(request)

	// Should fail due to foreign key constraint
	if err == nil {
		t.Error("Expected foreign key constraint error, got nil")
	}
}

func TestAddressRepository_GetAddresses_Integration(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestDBAfter(t, db)

	repo := addresses.NewAddressRepository(db)

	// Create test users and contacts
	user1 := users.User{Name: "User 1", Email: "user1@example.com", Password: "hash1"}
	user2 := users.User{Name: "User 2", Email: "user2@example.com", Password: "hash2"}
	db.Create(&user1)
	db.Create(&user2)

	contact1 := contacts.Contact{UserID: user1.ID, FirstName: "John", Email: "john@example.com"}
	contact2 := contacts.Contact{UserID: user2.ID, FirstName: "Jane", Email: "jane@example.com"}
	db.Create(&contact1)
	db.Create(&contact2)

	// Create addresses
	db.Create(&addresses.Address{ContactID: contact1.ID, Street: "123 Main St", City: "New York", Country: "USA"})
	db.Create(&addresses.Address{ContactID: contact1.ID, Street: "456 Oak Ave", City: "Boston", Country: "USA"})
	db.Create(&addresses.Address{ContactID: contact2.ID, Street: "789 Pine Rd", City: "Chicago", Country: "USA"})

	// Test get addresses for contact1 only
	result, err := repo.GetAddresses(contact1.ID, 1, 10, "")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// BUG CHECK: Should only return contact1's addresses (2), not all addresses (3)
	if len(result.Addresses) != 2 {
		t.Errorf("BUG FOUND: Expected 2 addresses for contact1, got %d (contact_id filter not working!)", len(result.Addresses))
	}

	if result.Total != 2 {
		t.Errorf("BUG FOUND: Expected total 2 for contact1, got %d (contact_id filter not working!)", result.Total)
	}

	// Verify no addresses from contact2 are returned
	for _, addr := range result.Addresses {
		if addr.ContactID != contact1.ID {
			t.Errorf("BUG FOUND: Address belongs to contact %d, not contact %d (authorization breach!)", addr.ContactID, contact1.ID)
		}
	}
}

func TestAddressRepository_GetAddresses_WithSearch_Integration(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestDBAfter(t, db)

	repo := addresses.NewAddressRepository(db)

	// Create test user and contact
	user1 := users.User{Name: "User 1", Email: "user1@example.com", Password: "hash1"}
	db.Create(&user1)

	contact1 := contacts.Contact{UserID: user1.ID, FirstName: "John", Email: "john@example.com"}
	db.Create(&contact1)

	// Create addresses
	db.Create(&addresses.Address{ContactID: contact1.ID, City: "New York", Country: "USA"})
	db.Create(&addresses.Address{ContactID: contact1.ID, City: "Boston", Country: "USA"})

	// Test search
	result, err := repo.GetAddresses(contact1.ID, 1, 10, "New York")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(result.Addresses) != 1 {
		t.Errorf("Expected 1 address, got %d", len(result.Addresses))
	}

	if result.Addresses[0].City != "New York" {
		t.Errorf("Expected city 'New York', got '%s'", result.Addresses[0].City)
	}
}

func TestAddressRepository_GetAddresses_Pagination_Integration(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestDBAfter(t, db)

	repo := addresses.NewAddressRepository(db)

	// Create test user and contact
	user1 := users.User{Name: "User 1", Email: "user1@example.com", Password: "hash1"}
	db.Create(&user1)

	contact1 := contacts.Contact{UserID: user1.ID, FirstName: "John", Email: "john@example.com"}
	db.Create(&contact1)

	// Create 5 addresses
	for i := 1; i <= 5; i++ {
		db.Create(&addresses.Address{ContactID: contact1.ID, City: "City", Country: "USA"})
	}

	// Test pagination
	result, err := repo.GetAddresses(contact1.ID, 1, 2, "")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(result.Addresses) != 2 {
		t.Errorf("Expected 2 addresses per page, got %d", len(result.Addresses))
	}

	if result.Total != 5 {
		t.Errorf("Expected total 5, got %d", result.Total)
	}

	if result.TotalPage != 3 {
		t.Errorf("Expected 3 total pages, got %d", result.TotalPage)
	}
}

func TestAddressRepository_FindAddressById_Integration(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestDBAfter(t, db)

	repo := addresses.NewAddressRepository(db)

	// Create test user and contact
	user1 := users.User{Name: "User 1", Email: "user1@example.com", Password: "hash1"}
	db.Create(&user1)

	contact1 := contacts.Contact{UserID: user1.ID, FirstName: "John", Email: "john@example.com"}
	db.Create(&contact1)

	// Create address
	createdAddress := addresses.Address{
		ContactID: contact1.ID,
		Street:    "123 Main St",
		City:      "New York",
		Country:   "USA",
	}
	db.Create(&createdAddress)

	// Find address
	address, err := repo.FindAddressById(createdAddress.ID)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if address.ID != createdAddress.ID {
		t.Errorf("Expected ID %d, got %d", createdAddress.ID, address.ID)
	}

	if address.City != "New York" {
		t.Errorf("Expected city 'New York', got '%s'", address.City)
	}
}

func TestAddressRepository_FindAddressById_WrongUser_Integration(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestDBAfter(t, db)

	repo := addresses.NewAddressRepository(db)

	// Create test users and contacts
	user1 := users.User{Name: "User 1", Email: "user1@example.com", Password: "hash1"}
	user2 := users.User{Name: "User 2", Email: "user2@example.com", Password: "hash2"}
	db.Create(&user1)
	db.Create(&user2)

	contact1 := contacts.Contact{UserID: user1.ID, FirstName: "John", Email: "john@example.com"}
	db.Create(&contact1)

	// Create address for user1's contact
	createdAddress := addresses.Address{
		ContactID: contact1.ID,
		City:      "New York",
		Country:   "USA",
	}
	db.Create(&createdAddress)

	// Try to find with user2's ID (should fail - authorization check)
	_, err := repo.FindAddressById(createdAddress.ID)

	if err == nil {
		t.Error("BUG FOUND: Should not find address belonging to different user (authorization breach!)")
	}
}

func TestAddressRepository_UpdateAddress_Integration(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestDBAfter(t, db)

	repo := addresses.NewAddressRepository(db)

	// Create test user and contact
	user1 := users.User{Name: "User 1", Email: "user1@example.com", Password: "hash1"}
	db.Create(&user1)

	contact1 := contacts.Contact{UserID: user1.ID, FirstName: "John", Email: "john@example.com"}
	db.Create(&contact1)

	// Create address
	createdAddress := addresses.Address{
		ContactID: contact1.ID,
		Street:    "123 Main St",
		City:      "New York",
		Country:   "USA",
	}
	db.Create(&createdAddress)

	// Update address
	updateAddress := &addresses.Address{
		Street: "456 Updated St",
		City:   "Boston",
	}

	_, err := repo.UpdateAddress(createdAddress.ID, updateAddress)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify update in database
	var address addresses.Address
	db.First(&address, createdAddress.ID)
	if address.Street != "456 Updated St" {
		t.Errorf("Expected street '456 Updated St', got '%s'", address.Street)
	}
}

func TestAddressRepository_UpdateAddress_PartialUpdate_Integration(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestDBAfter(t, db)

	repo := addresses.NewAddressRepository(db)

	// Create test user and contact
	user1 := users.User{Name: "User 1", Email: "user1@example.com", Password: "hash1"}
	db.Create(&user1)

	contact1 := contacts.Contact{UserID: user1.ID, FirstName: "John", Email: "john@example.com"}
	db.Create(&contact1)

	// Create address
	createdAddress := addresses.Address{
		ContactID:  contact1.ID,
		Street:     "123 Main St",
		City:       "New York",
		State:      "NY",
		PostalCode: "10001",
		Country:    "USA",
	}
	db.Create(&createdAddress)

	// Update only city
	updateAddress := &addresses.Address{
		City: "Boston",
	}

	_, err := repo.UpdateAddress(createdAddress.ID, updateAddress)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify update in database
	var address addresses.Address
	db.First(&address, createdAddress.ID)

	// BUG CHECK: Other fields should NOT be cleared
	if address.Street == "" {
		t.Error("BUG FOUND: Street was cleared when it should remain unchanged!")
	}
	if address.State == "" {
		t.Error("BUG FOUND: State was cleared when it should remain unchanged!")
	}
}

func TestAddressRepository_DeleteAddress_Integration(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestDBAfter(t, db)

	repo := addresses.NewAddressRepository(db)

	// Create test user and contact
	user1 := users.User{Name: "User 1", Email: "user1@example.com", Password: "hash1"}
	db.Create(&user1)

	contact1 := contacts.Contact{UserID: user1.ID, FirstName: "John", Email: "john@example.com"}
	db.Create(&contact1)

	// Create address
	createdAddress := addresses.Address{
		ContactID: contact1.ID,
		City:      "New York",
		Country:   "USA",
	}
	db.Create(&createdAddress)

	// Delete address
	err := repo.DeleteAddress(createdAddress.ID)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify deletion
	var address addresses.Address
	result := db.First(&address, createdAddress.ID)
	if result.Error == nil {
		t.Error("Address should be deleted but was still found")
	}
}

func TestAddressRepository_DeleteAddress_WrongUser_Integration(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestDBAfter(t, db)

	repo := addresses.NewAddressRepository(db)

	// Create test users and contacts
	user1 := users.User{Name: "User 1", Email: "user1@example.com", Password: "hash1"}
	user2 := users.User{Name: "User 2", Email: "user2@example.com", Password: "hash2"}
	db.Create(&user1)
	db.Create(&user2)

	contact1 := contacts.Contact{UserID: user1.ID, FirstName: "John", Email: "john@example.com"}
	db.Create(&contact1)

	// Create address for user1's contact
	createdAddress := addresses.Address{
		ContactID: contact1.ID,
		City:      "New York",
		Country:   "USA",
	}
	db.Create(&createdAddress)

	// Try to delete with user2's ID
	err := repo.DeleteAddress(createdAddress.ID)

	// Should succeed (no error) but address should NOT be deleted
	if err != nil {
		t.Logf("Delete returned error: %v", err)
	}

	// BUG CHECK: Verify address still exists (should not be deleted by wrong user)
	var address addresses.Address
	result := db.First(&address, createdAddress.ID)
	if result.Error != nil {
		t.Error("BUG FOUND: Address was deleted by wrong user (authorization breach!)")
	}
}
