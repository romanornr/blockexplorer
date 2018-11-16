package services

import (
	"github.com/romanornr/cyberchain/insightjson"
)

type addressDAO interface {
	// Get returns *insightjson.Address with the given hash.
	Get(addrID string) (*insightjson.AddressInfo, error)
	// Create saves the given address in the database.
	Create(addr *insightjson.AddressInfo) error
	// Delete deletes an address with the given hash from the database.
	Delete(addrID string) error
}

type AddressService struct {
	dao addressDAO
}

// NewBlockService creates a new BlockService with the given DAO.
func NewAddressService(dao addressDAO) *AddressService {
	return &AddressService{dao}
}

func (s *AddressService) Get(addrID string) (*insightjson.AddressInfo, error) {
	return s.dao.Get(addrID)
}

func (s *AddressService) Create(addr *insightjson.AddressInfo) error {
	if err := addr.Validate(); err != nil {
		return err
	}

	return s.dao.Create(addr)
}

func (s *AddressService) Delete(addrID string) error {
	_, err := s.dao.Get(addrID)
	if err != nil {
		return err
	}

	return s.dao.Delete(addrID)
}
