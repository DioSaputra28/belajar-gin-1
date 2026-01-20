package contacts

import (
	"errors"

	"gorm.io/gorm"
)

type ContactService interface {
	GetContacts(page, limit, user_id int, search string) (*GetContactsResponse, error)
	CreateContact(contact Contact) (*ContactResponse, error)
	FindContactById(id, user_id uint) (*Contact, error)
	UpdateContact(id, user_id uint, contact Contact) error
	DeleteContact(id, user_id uint) error
}

type contactService struct {
	repo ContactRepository
}

func NewContactService(repo ContactRepository) ContactService {
	return &contactService{repo: repo}
}

func (s *contactService) GetContacts(page, limit, user_id int, search string) (*GetContactsResponse, error) {
	response, err := s.repo.GetContacts(page, limit, user_id, search)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (s *contactService) CreateContact(contact Contact) (*ContactResponse, error) {
	contact_db, err := s.repo.CreateContact(contact)
	if err != nil {
		return nil, err
	}
	return contact_db, nil
}

func (s *contactService) FindContactById(id, user_id uint) (*Contact, error) {
	contact_db, err := s.repo.FindContactById(id, user_id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("contact not found")
		}
		return nil, err
	}
	return contact_db, nil
}

func (s *contactService) UpdateContact(id, user_id uint, contact Contact) error {
	contact_db, err := s.repo.FindContactById(id, user_id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("contact not found")
		}
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

	if err := s.repo.UpdateContact(id, user_id, contact_db); err != nil {
		return err
	}
	return nil
}

func (s *contactService) DeleteContact(id, user_id uint) error {
	_, err := s.repo.FindContactById(id, user_id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("contact not found")
		}
		return err
	}
	if err := s.repo.DeleteContact(id, user_id); err != nil {
		return err
	}
	return nil
}
