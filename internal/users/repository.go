package users

import "gorm.io/gorm"

type UserRepository interface {
	GetUsers(page, limit int, search string) ([]User, error)
	CreateUser(user CreateUserRequest) (*CreateUserResonse, error)
	UpdateUser(id uint, user UpdateUserRequest) error
	FindUserById(id uint) (*UserResponse, error)
	DeleteUser(id uint) error
}

type userRepository struct {
	db *gorm.DB	
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (u *userRepository) GetUsers(page, limit int, search string) ([]User, error) {
	var users []User
	if err := u.db.Where("name LIKE ? OR email LIKE ?", "%"+search+"%", "%"+search+"%").Offset((page - 1) * limit).Limit(limit).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (u *userRepository) CreateUser(user CreateUserRequest) (*CreateUserResonse, error) {
	if err := u.db.Create(&user).Error; err != nil {
		return nil, err
	}
	return &CreateUserResonse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (u *userRepository) UpdateUser(id uint, user UpdateUserRequest) error {
	var user_db User
	if err := u.db.Where("user_id = ?", id).First(&user_db).Error; err != nil {
		return err
	}

	user_db.Name = user.Name
	user_db.Email = user.Email

	if err := u.db.Save(&user_db).Error; err != nil {
		return err
	}
	return nil
}

func (u *userRepository) FindUserById(id uint) (*UserResponse, error) {
	var user User
	if err := u.db.Where("user_id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	
	return &UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (u *userRepository) DeleteUser(id uint) error {
	if err := u.db.Delete(&User{}, id).Error; err != nil {
		return err
	}
	return nil
}