package daos

import (
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/globalsign/mgo"
	"github.com/romanornr/cyberchain/insightjson"
)

type BlockDAO struct {
	session *mgo.Session
}

// NewBlockDAO creates a new BlockDAO
func NewBlockDAO(database dbConn) *BlockDAO {
	session, err := mgo.DialWithInfo(
		&mgo.DialInfo{
			Addrs:    dbHosts,
			Database: string(database),
			Timeout:  timeout10,
		},
	)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	return &BlockDAO{session}
}

func (dao *BlockDAO) Get(hash chainhash.Hash) (*insightjson.BlockResult, error) {
	panic("kek")
}

func (dao *BlockDAO) Create(block *insightjson.BlockResult) error {
	panic("kek")
}

func (dao *BlockDAO) Delete(hash chainhash.Hash) error {
	panic("kek")
}
