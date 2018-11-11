package services

import (
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/romanornr/cyberchain/insightjson"
)

type txDAO interface {
	// Get returns *insightjson.Tx with the given hash.
	Get(hash chainhash.Hash) (*insightjson.Tx, error)
	// Creates save the given transaction in the database.
	Create(tx *insightjson.Tx) error
	// Delete deletes a block with the given hash from the database.
	Delete(hash chainhash.Hash) error
}

type TxService struct {
	dao txDAO
}

// NewTxService creates a new BlockService with the given DAO.
func NewTxService(dao txDAO) *TxService {
	return &TxService{dao}
}

func (s *TxService) Get(hash chainhash.Hash) (*insightjson.Tx, error) {
	return s.dao.Get(hash)
}

func (s *TxService) Create(tx *insightjson.Tx) error {
	if err := tx.Validate(); err != nil {
		return err
	}

	// Here check if any of vin are not coinbase
	// if there are any, update existing tx

	return s.dao.Create(tx)
}

func (s *TxService) Delete(hash chainhash.Hash) error {
	_, err := s.dao.Get(hash)
	if err != nil {
		return err
	}

	return s.dao.Delete(hash)
}
