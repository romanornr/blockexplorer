package services

import (
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/romanornr/cyberchain/insightjson"
)

type addressDAO interface {
	// Get returns *insightjson.Address with the given hash.
	Get(hash *chainhash.Hash) (*insightjson.Address, error)
	// Create saves the given address in the database.
	Create(block *insightjson.Address) error
	// Delete deletes an address with the given hash from the database.
	Delete(hash *chainhash.Hash) error
}

type AddressService struct {
	dao addressDAO
}

// NewBlockService creates a new BlockService with the given DAO.
func NewAddressService(dao addressDAO) *AddressService {
	return &AddressService{dao}
}

func (s *AddressService) Get(hash *chainhash.Hash) (*insightjson.Address, error) {
	return s.dao.Get(hash)
}

func (s *AddressService) Create(addr *insightjson.Address) error {
	if err := addr.Validate(); err != nil {
		return err
	}

	return s.dao.Create(addr)
}

func (s *AddressService) Delete(hash *chainhash.Hash) error {
	_, err := s.dao.Get(hash)
	if err != nil {
		return err
	}

	return s.dao.Delete(hash)
}
