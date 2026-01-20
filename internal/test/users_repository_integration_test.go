package test

import (
	"testing"

	"github.com/DioSaputra28/belajar-gin-1/internal/users"
)

// ========== Users Repository Integration Tests ==========

func TestUserRepository_GetUsers_Integration(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestDBAfter(t, db)

	repo := users.NewUserRepository(db)

	// Create test users
	db.Create(&users.User{Name: "John Doe", Email: "john@example.com", Password: "hash1"})
	db.Create(&users.User{Name: "Jane Smith", Email: "jane@example.com", Password: "hash2"})
	db.Create(&users.User{Name: "Bob Wilson", Email: "bob@example.com", Password: "hash3"})

	// Test get users with pagination
	result, err := repo.GetUsers(1, 10, "")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(result.Data) != 3 {
		t.Errorf("Expected 3 users, got %d", len(result.Data))
	}

	if result.Total != 3 {
		t.Errorf("Expected total 3, got %d", result.Total)
	}
}

func TestUserRepository_GetUsers_WithSearch_Integration(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestDBAfter(t, db)

	repo := users.NewUserRepository(db)

	// Create test users
	db.Create(&users.User{Name: "John Doe", Email: "john@example.com", Password: "hash1"})
	db.Create(&users.User{Name: "Jane Smith", Email: "jane@example.com", Password: "hash2"})

	// Test search by name
	result, err := repo.GetUsers(1, 10, "john")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(result.Data) != 1 {
		t.Errorf("Expected 1 user, got %d", len(result.Data))
	}

	if result.Data[0].Name != "John Doe" {
		t.Errorf("Expected 'John Doe', got '%s'", result.Data[0].Name)
	}
}

func TestUserRepository_CreateUser_Integration(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestDBAfter(t, db)

	repo := users.NewUserRepository(db)

	request := users.CreateUserRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}

	response, err := repo.CreateUser(request)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if response.Name != "John Doe" {
		t.Errorf("Expected name 'John Doe', got '%s'", response.Name)
	}

	// Verify in database
	var user users.User
	db.Where("email = ?", "john@example.com").First(&user)
	if user.ID == 0 {
		t.Error("User was not saved to database")
	}
}

func TestUserRepository_UpdateUser_Integration(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestDBAfter(t, db)

	repo := users.NewUserRepository(db)

	// Create user
	createdUser := users.User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "hash1",
	}
	db.Create(&createdUser)

	// Update user
	updateRequest := users.UpdateUserRequest{
		Name:  "John Updated",
		Email: "john.updated@example.com",
	}

	err := repo.UpdateUser(createdUser.ID, updateRequest)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify update in database
	var user users.User
	db.First(&user, createdUser.ID)
	if user.Name != "John Updated" {
		t.Errorf("Expected name 'John Updated', got '%s'", user.Name)
	}
	if user.Email != "john.updated@example.com" {
		t.Errorf("Expected email 'john.updated@example.com', got '%s'", user.Email)
	}
}

func TestUserRepository_UpdateUser_PartialUpdate_Integration(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestDBAfter(t, db)

	repo := users.NewUserRepository(db)

	// Create user
	createdUser := users.User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "hash1",
	}
	db.Create(&createdUser)

	// Update only name (email should remain unchanged)
	updateRequest := users.UpdateUserRequest{
		Name: "John Updated",
	}

	err := repo.UpdateUser(createdUser.ID, updateRequest)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify update in database
	var user users.User
	db.First(&user, createdUser.ID)
	if user.Name != "John Updated" {
		t.Errorf("Expected name 'John Updated', got '%s'", user.Name)
	}
	// BUG CHECK: Email should NOT be empty
	if user.Email == "" {
		t.Error("BUG FOUND: Email was cleared when it should remain unchanged!")
	}
	if user.Email != "john@example.com" {
		t.Errorf("Expected email to remain 'john@example.com', got '%s'", user.Email)
	}
}

func TestUserRepository_FindUserById_Integration(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestDBAfter(t, db)

	repo := users.NewUserRepository(db)

	// Create user
	createdUser := users.User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "hash1",
	}
	db.Create(&createdUser)

	// Find user
	response, err := repo.FindUserById(createdUser.ID)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if response.ID != createdUser.ID {
		t.Errorf("Expected ID %d, got %d", createdUser.ID, response.ID)
	}
}

func TestUserRepository_DeleteUser_Integration(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestDBAfter(t, db)

	repo := users.NewUserRepository(db)

	// Create user
	createdUser := users.User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "hash1",
	}
	db.Create(&createdUser)

	// Delete user
	err := repo.DeleteUser(createdUser.ID)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify deletion (soft delete)
	var user users.User
	result := db.First(&user, createdUser.ID)
	if result.Error == nil {
		t.Error("User should be soft deleted but was still found")
	}
}
