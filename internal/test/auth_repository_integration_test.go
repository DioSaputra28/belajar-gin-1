package test

import (
	"testing"

	"github.com/DioSaputra28/belajar-gin-1/internal/auth"
	"github.com/DioSaputra28/belajar-gin-1/internal/common/utils"
	"github.com/DioSaputra28/belajar-gin-1/internal/users"
)

// ========== Auth Repository Integration Tests ==========

func TestAuthRepository_Register_Integration(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestDBAfter(t, db)

	repo := auth.NewAuthRepository(db)

	hashedPassword, _ := utils.HashPassword("password123")
	request := auth.RegisterRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: hashedPassword,
	}

	response, err := repo.Register(request)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if response.User.Name != "John Doe" {
		t.Errorf("Expected name 'John Doe', got '%s'", response.User.Name)
	}

	if response.User.Email != "john@example.com" {
		t.Errorf("Expected email 'john@example.com', got '%s'", response.User.Email)
	}

	// Verify user exists in database
	var user users.User
	db.Where("email = ?", "john@example.com").First(&user)
	if user.ID == 0 {
		t.Error("User was not saved to database")
	}
}

func TestAuthRepository_Login_Integration(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestDBAfter(t, db)

	repo := auth.NewAuthRepository(db)

	// Create user first
	hashedPassword, _ := utils.HashPassword("password123")
	db.Create(&users.User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: hashedPassword,
	})

	// Test login
	response, err := repo.Login("john@example.com")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if response.User.Email != "john@example.com" {
		t.Errorf("Expected email 'john@example.com', got '%s'", response.User.Email)
	}

	if response.AccessToken == "" {
		t.Error("Expected access token to be generated")
	}

	// Verify token was saved to database
	var user users.User
	db.Where("email = ?", "john@example.com").First(&user)
	if user.Token == "" {
		t.Error("Token was not saved to database")
	}
}

func TestAuthRepository_FindUserByEmail_Integration(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestDBAfter(t, db)

	repo := auth.NewAuthRepository(db)

	// Create user
	hashedPassword, _ := utils.HashPassword("password123")
	db.Create(&users.User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: hashedPassword,
	})

	// Test find
	user, err := repo.FindUserByEmail("john@example.com")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if user.Email != "john@example.com" {
		t.Errorf("Expected email 'john@example.com', got '%s'", user.Email)
	}
}

func TestAuthRepository_FindUserByToken_Integration(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestDBAfter(t, db)

	repo := auth.NewAuthRepository(db)

	// Create user with token
	hashedPassword, _ := utils.HashPassword("password123")
	testToken := "test-token-12345"
	db.Create(&users.User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: hashedPassword,
		Token:    testToken,
	})

	// Test find by token
	user, err := repo.FindUserByToken(testToken)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if user.Token != testToken {
		t.Errorf("Expected token '%s', got '%s'", testToken, user.Token)
	}
}

func TestAuthRepository_FindUserById_Integration(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestDBAfter(t, db)

	repo := auth.NewAuthRepository(db)

	// Create user
	hashedPassword, _ := utils.HashPassword("password123")
	createdUser := users.User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: hashedPassword,
	}
	db.Create(&createdUser)

	// Test find by ID
	user, err := repo.FindUserById(createdUser.ID)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if user.ID != createdUser.ID {
		t.Errorf("Expected ID %d, got %d", createdUser.ID, user.ID)
	}
}
