package services

import (
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/romanornr/cyberchain/insightjson"
)

type blockDAO interface {
	// Get returns *insightjson.BlockResult with the given hash.
	Get(hash chainhash.Hash) (*insightjson.BlockResult, error)
	// Create saves the given block in the database.
	Create(block *insightjson.BlockResult) error
	// Delete deletes a block with the given hash from the database.
	Delete(hash chainhash.Hash) error
}

type BlockService struct {
	dao blockDAO
}

// NewBlockService creates a new BlockService with the given DAO.
func NewBlockService(dao blockDAO) *BlockService {
	return &BlockService{dao}
}

func (s *BlockService) Get(hash chainhash.Hash) (*insightjson.BlockResult, error) {
	return s.dao.Get(hash)
}

func (s *BlockService) Create(block *insightjson.BlockResult) error {
	if err := block.Validate(); err != nil {
		return err
	}

	return s.dao.Create(block)
}

func (s *BlockService) Delete(hash chainhash.Hash) error {
	_, err := s.dao.Get(hash)
	if err != nil {
		return err
	}

	return s.dao.Delete(hash)
}
