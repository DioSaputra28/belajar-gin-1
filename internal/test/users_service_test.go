package test

import (
	"errors"
	"testing"
	"time"

	"github.com/DioSaputra28/belajar-gin-1/internal/users"
	"gorm.io/gorm"
)

// TestGetUsers_Success tests successful retrieval of users with pagination
func TestGetUsers_Success(t *testing.T) {
	mockRepo := &MockUserRepository{
		GetUsersFunc: func(page, limit int, search string) (*users.GetUsersResponse, error) {
			return &users.GetUsersResponse{
				Data: []users.User{
					{ID: 1, Name: "John Doe", Email: "john@example.com"},
					{ID: 2, Name: "Jane Smith", Email: "jane@example.com"},
				},
				Page:       page,
				Limit:      limit,
				Total:      2,
				TotalPages: 1,
			}, nil
		},
	}

	service := users.NewUserService(mockRepo)

	result, err := service.GetUsers(1, 10, "")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(result.Data) != 2 {
		t.Errorf("Expected 2 users, got %d", len(result.Data))
	}

	if result.Data[0].Name != "John Doe" {
		t.Errorf("Expected first user name 'John Doe', got '%s'", result.Data[0].Name)
	}

	if result.Total != 2 {
		t.Errorf("Expected total 2, got %d", result.Total)
	}

	if result.TotalPages != 1 {
		t.Errorf("Expected total pages 1, got %d", result.TotalPages)
	}
}

// TestGetUsers_WithSearch tests user retrieval with search query
func TestGetUsers_WithSearch(t *testing.T) {
	mockRepo := &MockUserRepository{
		GetUsersFunc: func(page, limit int, search string) (*users.GetUsersResponse, error) {
			if search == "john" {
				return &users.GetUsersResponse{
					Data: []users.User{
						{ID: 1, Name: "John Doe", Email: "john@example.com"},
					},
					Page:       page,
					Limit:      limit,
					Total:      1,
					TotalPages: 1,
				}, nil
			}
			return &users.GetUsersResponse{
				Data:       []users.User{},
				Page:       page,
				Limit:      limit,
				Total:      0,
				TotalPages: 0,
			}, nil
		},
	}

	service := users.NewUserService(mockRepo)

	result, err := service.GetUsers(1, 10, "john")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(result.Data) != 1 {
		t.Errorf("Expected 1 user, got %d", len(result.Data))
	}

	if result.Data[0].Email != "john@example.com" {
		t.Errorf("Expected email 'john@example.com', got '%s'", result.Data[0].Email)
	}
}

// TestGetUsers_EmptyResult tests user retrieval with no results
func TestGetUsers_EmptyResult(t *testing.T) {
	mockRepo := &MockUserRepository{
		GetUsersFunc: func(page, limit int, search string) (*users.GetUsersResponse, error) {
			return &users.GetUsersResponse{
				Data:       []users.User{},
				Page:       page,
				Limit:      limit,
				Total:      0,
				TotalPages: 0,
			}, nil
		},
	}

	service := users.NewUserService(mockRepo)

	result, err := service.GetUsers(1, 10, "nonexistent")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(result.Data) != 0 {
		t.Errorf("Expected 0 users, got %d", len(result.Data))
	}
}

// TestGetUsers_DatabaseError tests user retrieval with database error
func TestGetUsers_DatabaseError(t *testing.T) {
	mockRepo := &MockUserRepository{
		GetUsersFunc: func(page, limit int, search string) (*users.GetUsersResponse, error) {
			return nil, errors.New("database connection error")
		},
	}

	service := users.NewUserService(mockRepo)

	_, err := service.GetUsers(1, 10, "")

	if err == nil {
		t.Error("Expected database error, got nil")
	}
}

// TestCreateUser_Success tests successful user creation
func TestCreateUser_Success(t *testing.T) {
	mockRepo := &MockUserRepository{
		CreateUserFunc: func(user users.CreateUserRequest) (*users.CreateUserResonse, error) {
			return &users.CreateUserResonse{
				ID:    1,
				Name:  user.Name,
				Email: user.Email,
			}, nil
		},
	}

	service := users.NewUserService(mockRepo)

	request := users.CreateUserRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}

	response, err := service.CreateUser(request)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if response.Name != "John Doe" {
		t.Errorf("Expected name 'John Doe', got '%s'", response.Name)
	}

	if response.Email != "john@example.com" {
		t.Errorf("Expected email 'john@example.com', got '%s'", response.Email)
	}
}

// TestCreateUser_MinimumData tests user creation with minimum valid data
func TestCreateUser_MinimumData(t *testing.T) {
	mockRepo := &MockUserRepository{
		CreateUserFunc: func(user users.CreateUserRequest) (*users.CreateUserResonse, error) {
			return &users.CreateUserResonse{
				ID:    1,
				Name:  user.Name,
				Email: user.Email,
			}, nil
		},
	}

	service := users.NewUserService(mockRepo)

	request := users.CreateUserRequest{
		Name:     "Joe",
		Email:    "joe@example.com",
		Password: "pass12",
	}

	response, err := service.CreateUser(request)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if response.Name != "Joe" {
		t.Errorf("Expected name 'Joe', got '%s'", response.Name)
	}
}

// TestCreateUser_DatabaseError tests user creation with database error
func TestCreateUser_DatabaseError(t *testing.T) {
	mockRepo := &MockUserRepository{
		CreateUserFunc: func(user users.CreateUserRequest) (*users.CreateUserResonse, error) {
			return nil, errors.New("duplicate key value violates unique constraint")
		},
	}

	service := users.NewUserService(mockRepo)

	request := users.CreateUserRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}

	_, err := service.CreateUser(request)

	if err == nil {
		t.Error("Expected database error, got nil")
	}
}

// TestUpdateUser_Success tests successful user update
func TestUpdateUser_Success(t *testing.T) {
	mockRepo := &MockUserRepository{
		FindUserByIdFunc: func(id uint) (*users.UserResponse, error) {
			return &users.UserResponse{
				ID:    id,
				Name:  "John Doe",
				Email: "john@example.com",
			}, nil
		},
		UpdateUserFunc: func(id uint, user users.UpdateUserRequest) error {
			return nil
		},
	}

	service := users.NewUserService(mockRepo)

	request := users.UpdateUserRequest{
		Name:  "John Updated",
		Email: "john.updated@example.com",
	}

	err := service.UpdateUser(1, request)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

// TestUpdateUser_NameOnly tests updating only the name
func TestUpdateUser_NameOnly(t *testing.T) {
	mockRepo := &MockUserRepository{
		FindUserByIdFunc: func(id uint) (*users.UserResponse, error) {
			return &users.UserResponse{
				ID:    id,
				Name:  "John Doe",
				Email: "john@example.com",
			}, nil
		},
		UpdateUserFunc: func(id uint, user users.UpdateUserRequest) error {
			return nil
		},
	}

	service := users.NewUserService(mockRepo)

	request := users.UpdateUserRequest{
		Name: "John Updated",
	}

	err := service.UpdateUser(1, request)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

// TestUpdateUser_EmailOnly tests updating only the email
func TestUpdateUser_EmailOnly(t *testing.T) {
	mockRepo := &MockUserRepository{
		FindUserByIdFunc: func(id uint) (*users.UserResponse, error) {
			return &users.UserResponse{
				ID:    id,
				Name:  "John Doe",
				Email: "john@example.com",
			}, nil
		},
		UpdateUserFunc: func(id uint, user users.UpdateUserRequest) error {
			return nil
		},
	}

	service := users.NewUserService(mockRepo)

	request := users.UpdateUserRequest{
		Email: "john.new@example.com",
	}

	err := service.UpdateUser(1, request)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

// TestUpdateUser_NotFound tests updating non-existent user
func TestUpdateUser_NotFound(t *testing.T) {
	mockRepo := &MockUserRepository{
		FindUserByIdFunc: func(id uint) (*users.UserResponse, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}

	service := users.NewUserService(mockRepo)

	request := users.UpdateUserRequest{
		Name: "John Updated",
	}

	err := service.UpdateUser(999, request)

	if err == nil {
		t.Error("Expected error for non-existent user, got nil")
	}

	if err.Error() != "user not found" {
		t.Errorf("Expected 'user not found' error, got '%s'", err.Error())
	}
}

// TestUpdateUser_DatabaseError tests update with database error
func TestUpdateUser_DatabaseError(t *testing.T) {
	mockRepo := &MockUserRepository{
		FindUserByIdFunc: func(id uint) (*users.UserResponse, error) {
			return &users.UserResponse{
				ID:    id,
				Name:  "John Doe",
				Email: "john@example.com",
			}, nil
		},
		UpdateUserFunc: func(id uint, user users.UpdateUserRequest) error {
			return errors.New("database connection error")
		},
	}

	service := users.NewUserService(mockRepo)

	request := users.UpdateUserRequest{
		Name: "John Updated",
	}

	err := service.UpdateUser(1, request)

	if err == nil {
		t.Error("Expected database error, got nil")
	}
}

// TestFindUserById_Success tests successful user retrieval by ID
func TestFindUserById_Success(t *testing.T) {
	now := time.Now()
	mockRepo := &MockUserRepository{
		FindUserByIdFunc: func(id uint) (*users.UserResponse, error) {
			return &users.UserResponse{
				ID:        id,
				Name:      "John Doe",
				Email:     "john@example.com",
				CreatedAt: now,
				UpdatedAt: now,
			}, nil
		},
	}

	service := users.NewUserService(mockRepo)

	response, err := service.FindUserById(1)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if response.ID != 1 {
		t.Errorf("Expected user ID 1, got %d", response.ID)
	}

	if response.Name != "John Doe" {
		t.Errorf("Expected name 'John Doe', got '%s'", response.Name)
	}
}

// TestFindUserById_NotFound tests finding non-existent user
func TestFindUserById_NotFound(t *testing.T) {
	mockRepo := &MockUserRepository{
		FindUserByIdFunc: func(id uint) (*users.UserResponse, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}

	service := users.NewUserService(mockRepo)

	_, err := service.FindUserById(999)

	if err == nil {
		t.Error("Expected error for non-existent user, got nil")
	}

	if err.Error() != "user not found" {
		t.Errorf("Expected 'user not found' error, got '%s'", err.Error())
	}
}

// TestFindUserById_DatabaseError tests finding user with database error
func TestFindUserById_DatabaseError(t *testing.T) {
	mockRepo := &MockUserRepository{
		FindUserByIdFunc: func(id uint) (*users.UserResponse, error) {
			return nil, errors.New("database connection error")
		},
	}

	service := users.NewUserService(mockRepo)

	_, err := service.FindUserById(1)

	if err == nil {
		t.Error("Expected database error, got nil")
	}
}

// TestDeleteUser_Success tests successful user deletion
func TestDeleteUser_Success(t *testing.T) {
	mockRepo := &MockUserRepository{
		FindUserByIdFunc: func(id uint) (*users.UserResponse, error) {
			return &users.UserResponse{
				ID:    id,
				Name:  "John Doe",
				Email: "john@example.com",
			}, nil
		},
		DeleteUserFunc: func(id uint) error {
			return nil
		},
	}

	service := users.NewUserService(mockRepo)

	err := service.DeleteUser(1)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

// TestDeleteUser_NotFound tests deleting non-existent user
func TestDeleteUser_NotFound(t *testing.T) {
	mockRepo := &MockUserRepository{
		FindUserByIdFunc: func(id uint) (*users.UserResponse, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}

	service := users.NewUserService(mockRepo)

	err := service.DeleteUser(999)

	if err == nil {
		t.Error("Expected error for non-existent user, got nil")
	}

	if err.Error() != "user not found" {
		t.Errorf("Expected 'user not found' error, got '%s'", err.Error())
	}
}

// TestDeleteUser_DatabaseError tests deletion with database error
func TestDeleteUser_DatabaseError(t *testing.T) {
	mockRepo := &MockUserRepository{
		FindUserByIdFunc: func(id uint) (*users.UserResponse, error) {
			return &users.UserResponse{
				ID:    id,
				Name:  "John Doe",
				Email: "john@example.com",
			}, nil
		},
		DeleteUserFunc: func(id uint) error {
			return errors.New("database connection error")
		},
	}

	service := users.NewUserService(mockRepo)

	err := service.DeleteUser(1)

	if err == nil {
		t.Error("Expected database error, got nil")
	}
}
