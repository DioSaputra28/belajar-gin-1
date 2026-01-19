package auth

import (
	"github.com/DioSaputra28/belajar-gin-1/internal/common/utils"
	"github.com/DioSaputra28/belajar-gin-1/internal/users"
	"gorm.io/gorm"
)

type AuthRepository interface {
	Register(request RegisterRequest) (AuthResponse, error)
	Login(email string) (AuthResponse, error)
	FindUserByEmail(email string) (*users.User, error)
	FindUserByToken(token string) (*users.User, error)
	FindUserById(id uint) (*users.User, error)
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}

func (a *authRepository) Register(request RegisterRequest) (AuthResponse, error) {
	user := users.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: request.Password,
	}

	result := a.db.Create(&user)
	if result.Error != nil {
		return AuthResponse{}, result.Error
	}

	return AuthResponse{
		User:         UserData{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
		AccessToken:  "",
		RefreshToken: "",
	}, nil
}

func (a *authRepository) Login(email string) (AuthResponse, error) {
	token := utils.GenerateToken()

	user, err := a.FindUserByEmail(email)
	if err != nil {
		return AuthResponse{}, err
	}

	user.Token = token

	result := a.db.Save(&user)
	if result.Error != nil {
		return AuthResponse{}, result.Error
	}

	return AuthResponse{
		User:         UserData{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
		AccessToken:  token,
		RefreshToken: "",
	}, nil
}

func (a *authRepository) FindUserByEmail(email string) (*users.User, error) {
	var user users.User
	if err := a.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (a *authRepository) FindUserByToken(token string) (*users.User, error) {
	var user users.User
	if err := a.db.Where("token = ?", token).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (a *authRepository) FindUserById(id uint) (*users.User, error) {
	var user users.User
	if err := a.db.Where("user_id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
