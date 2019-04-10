package services

import (
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/romanornr/blockexplorer/insightjson"
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

func (b *BlockService) Get(hash chainhash.Hash) (*insightjson.BlockResult, error) {
	return b.dao.Get(hash)
}

func (b *BlockService) Create(block *insightjson.BlockResult) error {
	if err := block.Validate(); err != nil {
		return err
	}
	if err := b.dao.Create(block); err != nil {
		return err
	}

	return nil
}

func (b *BlockService) Delete(hash chainhash.Hash) error {
	_, err := b.dao.Get(hash)
	if err != nil {
		return err
	}
	err = b.dao.Delete(hash)

	return err
}
