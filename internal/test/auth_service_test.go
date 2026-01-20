package test

import (
	"errors"
	"testing"

	"github.com/DioSaputra28/belajar-gin-1/internal/auth"
	"github.com/DioSaputra28/belajar-gin-1/internal/common/utils"
	"github.com/DioSaputra28/belajar-gin-1/internal/users"
	"gorm.io/gorm"
)

// TestRegister_Success tests successful user registration
func TestRegister_Success(t *testing.T) {
	mockRepo := &MockAuthRepository{
		FindUserByEmailFunc: func(email string) (*users.User, error) {
			return nil, gorm.ErrRecordNotFound
		},
		RegisterFunc: func(request auth.RegisterRequest) (auth.AuthResponse, error) {
			return auth.AuthResponse{
				User: auth.UserData{
					ID:    1,
					Name:  request.Name,
					Email: request.Email,
				},
				AccessToken:  "",
				RefreshToken: "",
			}, nil
		},
	}

	service := auth.NewAuthService(mockRepo)

	request := auth.RegisterRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}

	response, err := service.Register(request)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if response.User.Name != "John Doe" {
		t.Errorf("Expected name 'John Doe', got '%s'", response.User.Name)
	}

	if response.User.Email != "john@example.com" {
		t.Errorf("Expected email 'john@example.com', got '%s'", response.User.Email)
	}
}

// TestRegister_DuplicateEmail tests registration with existing email
func TestRegister_DuplicateEmail(t *testing.T) {
	existingUser := &users.User{
		ID:    1,
		Name:  "Existing User",
		Email: "john@example.com",
	}

	mockRepo := &MockAuthRepository{
		FindUserByEmailFunc: func(email string) (*users.User, error) {
			return existingUser, nil
		},
	}

	service := auth.NewAuthService(mockRepo)

	request := auth.RegisterRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}

	_, err := service.Register(request)

	if err == nil {
		t.Error("Expected error for duplicate email, got nil")
	}

	if err.Error() != "user already exists" {
		t.Errorf("Expected 'user already exists' error, got '%s'", err.Error())
	}
}

// TestRegister_DatabaseError tests registration with database error
func TestRegister_DatabaseError(t *testing.T) {
	mockRepo := &MockAuthRepository{
		FindUserByEmailFunc: func(email string) (*users.User, error) {
			return nil, errors.New("database connection error")
		},
	}

	service := auth.NewAuthService(mockRepo)

	request := auth.RegisterRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}

	_, err := service.Register(request)

	if err == nil {
		t.Error("Expected database error, got nil")
	}
}

// TestRegister_RepositoryError tests registration failure at repository level
func TestRegister_RepositoryError(t *testing.T) {
	mockRepo := &MockAuthRepository{
		FindUserByEmailFunc: func(email string) (*users.User, error) {
			return nil, gorm.ErrRecordNotFound
		},
		RegisterFunc: func(request auth.RegisterRequest) (auth.AuthResponse, error) {
			return auth.AuthResponse{}, errors.New("failed to create user")
		},
	}

	service := auth.NewAuthService(mockRepo)

	request := auth.RegisterRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}

	_, err := service.Register(request)

	if err == nil {
		t.Error("Expected repository error, got nil")
	}

	if err.Error() != "failed to create user" {
		t.Errorf("Expected 'failed to create user' error, got '%s'", err.Error())
	}
}

// TestLogin_Success tests successful login
func TestLogin_Success(t *testing.T) {
	// Generate a proper bcrypt hash for "password123" using the project's utility
	hashedPassword, err := utils.HashPassword("password123")
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	mockRepo := &MockAuthRepository{
		FindUserByEmailFunc: func(email string) (*users.User, error) {
			return &users.User{
				ID:       1,
				Name:     "John Doe",
				Email:    email,
				Password: hashedPassword,
			}, nil
		},
		LoginFunc: func(email string) (auth.AuthResponse, error) {
			return auth.AuthResponse{
				User: auth.UserData{
					ID:    1,
					Name:  "John Doe",
					Email: email,
				},
				AccessToken:  "test-token-123",
				RefreshToken: "",
			}, nil
		},
	}

	service := auth.NewAuthService(mockRepo)

	response, err := service.Login("john@example.com", "password123")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if response.AccessToken != "test-token-123" {
		t.Errorf("Expected token 'test-token-123', got '%s'", response.AccessToken)
	}
}

// TestLogin_UserNotFound tests login with non-existent user
func TestLogin_UserNotFound(t *testing.T) {
	mockRepo := &MockAuthRepository{
		FindUserByEmailFunc: func(email string) (*users.User, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}

	service := auth.NewAuthService(mockRepo)

	_, err := service.Login("nonexistent@example.com", "password123")

	if err == nil {
		t.Error("Expected error for non-existent user, got nil")
	}
}

// TestLogin_IncorrectPassword tests login with wrong password
func TestLogin_IncorrectPassword(t *testing.T) {
	hashedPassword, err := utils.HashPassword("password123")
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	mockRepo := &MockAuthRepository{
		FindUserByEmailFunc: func(email string) (*users.User, error) {
			return &users.User{
				ID:       1,
				Name:     "John Doe",
				Email:    email,
				Password: hashedPassword,
			}, nil
		},
	}

	service := auth.NewAuthService(mockRepo)

	_, err = service.Login("john@example.com", "wrongpassword")

	if err == nil {
		t.Error("Expected error for incorrect password, got nil")
	}

	if err.Error() != "email or password is incorrect" {
		t.Errorf("Expected 'email or password is incorrect' error, got '%s'", err.Error())
	}
}

// TestLogin_DatabaseError tests login with database error
func TestLogin_DatabaseError(t *testing.T) {
	mockRepo := &MockAuthRepository{
		FindUserByEmailFunc: func(email string) (*users.User, error) {
			return nil, errors.New("database connection error")
		},
	}

	service := auth.NewAuthService(mockRepo)

	_, err := service.Login("john@example.com", "password123")

	if err == nil {
		t.Error("Expected database error, got nil")
	}
}

// TestMe_Success tests successful user retrieval
func TestMe_Success(t *testing.T) {
	mockRepo := &MockAuthRepository{
		FindUserByIdFunc: func(id uint) (*users.User, error) {
			return &users.User{
				ID:    id,
				Name:  "John Doe",
				Email: "john@example.com",
				Token: "test-token-123",
			}, nil
		},
	}

	service := auth.NewAuthService(mockRepo)

	response, err := service.Me(1)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if response.User.ID != 1 {
		t.Errorf("Expected user ID 1, got %d", response.User.ID)
	}

	if response.User.Name != "John Doe" {
		t.Errorf("Expected name 'John Doe', got '%s'", response.User.Name)
	}

	if response.AccessToken != "test-token-123" {
		t.Errorf("Expected token 'test-token-123', got '%s'", response.AccessToken)
	}
}

// TestMe_UserNotFound tests Me with non-existent user
func TestMe_UserNotFound(t *testing.T) {
	mockRepo := &MockAuthRepository{
		FindUserByIdFunc: func(id uint) (*users.User, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}

	service := auth.NewAuthService(mockRepo)

	_, err := service.Me(999)

	if err == nil {
		t.Error("Expected error for non-existent user, got nil")
	}

	if err.Error() != "user not found" {
		t.Errorf("Expected 'user not found' error, got '%s'", err.Error())
	}
}

// TestMe_DatabaseError tests Me with database error (non-RecordNotFound)
func TestMe_DatabaseError(t *testing.T) {
	mockRepo := &MockAuthRepository{
		FindUserByIdFunc: func(id uint) (*users.User, error) {
			// Return database error (not RecordNotFound)
			return nil, errors.New("database connection error")
		},
	}

	service := auth.NewAuthService(mockRepo)

	// After the bug fix, non-RecordNotFound errors should be properly propagated
	_, err := service.Me(1)

	if err == nil {
		t.Error("Expected database error, got nil")
	}

	if err.Error() != "database connection error" {
		t.Errorf("Expected 'database connection error', got '%s'", err.Error())
	}
}
