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
	return &BlockService{}
}

func (b *BlockService) Get(hash chainhash.Hash) (*insightjson.BlockResult, error) {
	return b.dao.Get(hash)
}

func (b *BlockService) Create(block *insightjson.BlockResult) error {
	panic("implement me")
}

func (b *BlockService) Delete(hash chainhash.Hash) error {
	panic("implement me")
}
