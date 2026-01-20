package test

import (
	"github.com/DioSaputra28/belajar-gin-1/internal/auth"
	"github.com/DioSaputra28/belajar-gin-1/internal/users"
)

// MockAuthRepository implements auth.AuthRepository interface
type MockAuthRepository struct {
	RegisterFunc        func(request auth.RegisterRequest) (auth.AuthResponse, error)
	LoginFunc           func(email string) (auth.AuthResponse, error)
	FindUserByEmailFunc func(email string) (*users.User, error)
	FindUserByTokenFunc func(token string) (*users.User, error)
	FindUserByIdFunc    func(id uint) (*users.User, error)
}

// Register implements auth.AuthRepository
func (m *MockAuthRepository) Register(request auth.RegisterRequest) (auth.AuthResponse, error) {
	if m.RegisterFunc != nil {
		return m.RegisterFunc(request)
	}
	return auth.AuthResponse{}, nil
}

// Login implements auth.AuthRepository
func (m *MockAuthRepository) Login(email string) (auth.AuthResponse, error) {
	if m.LoginFunc != nil {
		return m.LoginFunc(email)
	}
	return auth.AuthResponse{}, nil
}

// FindUserByEmail implements auth.AuthRepository
func (m *MockAuthRepository) FindUserByEmail(email string) (*users.User, error) {
	if m.FindUserByEmailFunc != nil {
		return m.FindUserByEmailFunc(email)
	}
	return nil, nil
}

// FindUserByToken implements auth.AuthRepository
func (m *MockAuthRepository) FindUserByToken(token string) (*users.User, error) {
	if m.FindUserByTokenFunc != nil {
		return m.FindUserByTokenFunc(token)
	}
	return nil, nil
}

// FindUserById implements auth.AuthRepository
func (m *MockAuthRepository) FindUserById(id uint) (*users.User, error) {
	if m.FindUserByIdFunc != nil {
		return m.FindUserByIdFunc(id)
	}
	return nil, nil
}
