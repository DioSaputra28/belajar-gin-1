package users

import "gorm.io/gorm"

type UserRepository interface {
	GetUsers(page, limit int, search string) (*GetUsersResponse, error)
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

func (u *userRepository) GetUsers(page, limit int, search string) (*GetUsersResponse, error) {
	var users []User
	var total int64

	// Build query dengan search filter
	query := u.db.Model(&User{})
	if search != "" {
		query = query.Where("name LIKE ? OR email LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// Count total records (sebelum pagination)
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// Get paginated data
	if err := query.Offset((page - 1) * limit).Limit(limit).Find(&users).Error; err != nil {
		return nil, err
	}

	// Hitung total pages
	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++ // Tambah 1 jika ada sisa
	}

	return &GetUsersResponse{
		Data:       users,
		Page:       page,
		Limit:      limit,
		Total:      int(total),
		TotalPages: totalPages,
	}, nil
}

func (u *userRepository) CreateUser(user CreateUserRequest) (*CreateUserResonse, error) {
	var userModel User
	userModel.Name = user.Name
	userModel.Email = user.Email
	userModel.Password = user.Password
	if err := u.db.Create(&userModel).Error; err != nil {
		return nil, err
	}
	return &CreateUserResonse{
		ID:    userModel.ID,
		Name:  userModel.Name,
		Email: userModel.Email,
	}, nil
}

func (u *userRepository) UpdateUser(id uint, user UpdateUserRequest) error {
	var user_db User
	if err := u.db.Where("user_id = ?", id).First(&user_db).Error; err != nil {
		return err
	}

	if user.Name != "" {
		user_db.Name = user.Name
	}
	if user.Email != "" {
		user_db.Email = user.Email
	}

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
