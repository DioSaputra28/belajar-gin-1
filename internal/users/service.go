package users

import (
	"errors"

	"gorm.io/gorm"
)

type UserService interface {
	GetUsers(page, limit int, search string) (*GetUsersResponse, error)
	CreateUser(user CreateUserRequest) (*CreateUserResonse, error)
	UpdateUser(id uint, user UpdateUserRequest) error
	FindUserById(id uint) (*UserResponse, error)
	DeleteUser(id uint) error
}

type userService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) GetUsers(page, limit int, search string) (*GetUsersResponse, error) {
	response, err := s.repo.GetUsers(page, limit, search)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *userService) CreateUser(user CreateUserRequest) (*CreateUserResonse, error) {
	user_db, err := s.repo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return user_db, nil
}

func (s *userService) UpdateUser(id uint, user UpdateUserRequest) error {

	_, err := s.repo.FindUserById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	err = s.repo.UpdateUser(id, user)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	return nil
}

func (s *userService) FindUserById(id uint) (*UserResponse, error) {
	user_db, err := s.repo.FindUserById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user_db, nil
}

func (s *userService) DeleteUser(id uint) error {
	_, err := s.repo.FindUserById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	err = s.repo.DeleteUser(id)
	if err != nil {
		return err
	}

	return nil
}
