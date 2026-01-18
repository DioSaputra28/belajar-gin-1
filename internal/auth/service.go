package auth

import (
	"errors"

	"github.com/DioSaputra28/belajar-gin-1/internal/common/utils"
)


type AuthService interface {
	Register(request RegisterRequest) (AuthResponse, error)
	Login(email string, password string) (AuthResponse, error)
}

type authService struct {
	repo AuthRepository
}

func NewAuthService(repo AuthRepository) AuthService {
	return &authService{repo: repo}
}

func (s *authService) Register(request RegisterRequest) (AuthResponse, error) {
	user, err := s.repo.FindUserByEmail(request.Email)
	if err != nil {
		return AuthResponse{}, err
	}

	if user != nil {
		return AuthResponse{}, errors.New("user already exists")
	}

	hashedPassword, err := utils.HashPassword(request.Password)
	if err != nil {
		return AuthResponse{}, err
	}

	request.Password = hashedPassword

	result, err := s.repo.Register(request)
	if err != nil {
		return AuthResponse{}, err
	}

	return result, nil
}

func (s *authService) Login(email string, password string) (AuthResponse, error) {
	user, err := s.repo.FindUserByEmail(email)
	if err != nil {
		return AuthResponse{}, err
	}

	if user == nil {
		return AuthResponse{}, errors.New("user not found")
	}

	err = utils.CheckPassword(password, user.Password)
	if err != nil {
		return AuthResponse{}, errors.New("email or password is incorrect")
	}


	result, err := s.repo.Login(email)
	if err != nil {
		return AuthResponse{}, err
	}

	return result, nil
}