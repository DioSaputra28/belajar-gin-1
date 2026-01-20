package addresses

import (
	"errors"

	"gorm.io/gorm"
)

type AddressService interface {
	CreateAddress(user_id uint, address CreateAddressRequest) (*AddressResponse, error)
	GetAddresses(user_id, contact_id uint, page int, limit int, search string) (*GetAddressesResponse, error)
	UpdateAddress(user_id, address_id uint, address UpdateAddressRequest) (*AddressResponse, error)
	FindAddressById(user_id, address_id uint) (*Address, error)
	DeleteAddress(user_id, address_id uint) error
}

type addressService struct {
	repo AddressRepository
}

func NewAddressService(repo AddressRepository) AddressService {
	return &addressService{repo: repo}
}

func (s *addressService) CreateAddress(user_id uint, address CreateAddressRequest) (*AddressResponse, error) {
	contact_db, err := s.repo.FindContactById(address.ContactID, user_id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("contact not found")
		}
		return nil, err
	}

	if contact_db == nil {
		return nil, errors.New("contact not found")
	}

	result, err := s.repo.CreateAddress(address)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *addressService) GetAddresses(user_id, contact_id uint, page int, limit int, search string) (*GetAddressesResponse, error) {
	_, err := s.repo.FindContactById(contact_id, user_id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("contact not found")
		}
		return nil, err
	}
	
	address, err := s.repo.GetAddresses(contact_id, page, limit, search)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("contact not found")
		}
		return nil, err
	}

	return address, nil
}

func (s *addressService) UpdateAddress(user_id, address_id uint, address UpdateAddressRequest) (*AddressResponse, error) {
	address_db, err := s.repo.FindAddressById(address_id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("address not found")
		}
		return nil, err
	}

	if address.City != "" {
		address_db.City = address.City
	}

	if address.Country != "" {
		address_db.Country = address.Country
	}

	if address.PostalCode != "" {
		address_db.PostalCode = address.PostalCode
	}

	if address.State != "" {
		address_db.State = address.State
	}

	if address.Street != "" {
		address_db.Street = address.Street
	}

	result, err := s.repo.UpdateAddress(address_id, address_db)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("address not found")
		}
		return nil, err
	}

	return result, nil
}

func (s *addressService) FindAddressById(user_id, address_id uint) (*Address, error) {
	result, err := s.repo.FindAddressById(address_id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("address not found")
		}
		return nil, err
	}

	return result, nil
}

func (s *addressService) DeleteAddress(user_id, address_id uint) error {
	_, err := s.repo.FindAddressById(address_id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("address not found")
		}
		return err
	}

	err = s.repo.DeleteAddress(address_id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("contact not found")
		}
		return err
	}

	return nil
}