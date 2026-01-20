package contacts

import "gorm.io/gorm"

type ContactRepository interface {
	GetContacts(page, limit, user_id int, search string) (*GetContactsResponse, error)
	CreateContact(contact Contact) (*ContactResponse, error)
	FindContactById(id, user_id uint) (*Contact, error)
	UpdateContact(id, user_id uint, contact *Contact) error
	DeleteContact(id, user_id uint) error
}

type contactRepository struct {
	db *gorm.DB
}

func NewContactRepository(db *gorm.DB) ContactRepository {
	return &contactRepository{db: db}
}

func (c *contactRepository) GetContacts(page, limit, user_id int, search string) (*GetContactsResponse, error) {
	var contacts []Contact
	var total int64

	// Build query dengan search filter
	query := c.db.Model(&Contact{}).Where("user_id = ?", user_id)
	if search != "" {
		query = query.Where("first_name LIKE ? OR last_name LIKE ? OR email LIKE ? OR phone LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// Count total records (sebelum pagination)
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// Get paginated data
	if err := query.Offset((page - 1) * limit).Limit(limit).Find(&contacts).Error; err != nil {
		return nil, err
	}

	// Hitung total pages
	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++ // Tambah 1 jika ada sisa
	}

	return &GetContactsResponse{
		Data:       contacts,
		Page:       page,
		Limit:      limit,
		Total:      int(total),
		TotalPages: totalPages,
	}, nil
}

func (c *contactRepository) CreateContact(contact Contact) (*ContactResponse, error) {
	if err := c.db.Create(&contact).Error; err != nil {
		return nil, err
	}
	return &ContactResponse{
		ID:    contact.ID,
		FirstName:  contact.FirstName,
		LastName: contact.LastName,
		Email: contact.Email,
		Phone: contact.Phone,
	}, nil
}

func (c *contactRepository) FindContactById(id, user_id uint) (*Contact, error) {
	var contact Contact
	if err := c.db.Where("contact_id = ? AND user_id = ?", id, user_id).First(&contact).Error; err != nil {
		return nil, err
	}
	return &contact, nil
}

func (c *contactRepository) UpdateContact(id uint, user_id uint, contact *Contact) error {
	var contact_db Contact
	if err := c.db.Where("contact_id = ? AND user_id = ?", id, user_id).First(&contact_db).Error; err != nil {
		return err
	}

	if contact.FirstName != "" {
		contact_db.FirstName = contact.FirstName
	}
	if contact.LastName != "" {
		contact_db.LastName = contact.LastName
	}
	if contact.Email != "" {
		contact_db.Email = contact.Email
	}
	if contact.Phone != "" {
		contact_db.Phone = contact.Phone
	}

	if err := c.db.Save(&contact_db).Error; err != nil {
		return err
	}
	return nil
}

func (c *contactRepository) DeleteContact(id uint, user_id uint) error {
	if err := c.db.Where("contact_id = ? AND user_id = ?", id, user_id).Delete(&Contact{}).Error; err != nil {
		return err
	}
	return nil
}
