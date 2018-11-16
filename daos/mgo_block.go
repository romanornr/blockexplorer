package daos

import (
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/globalsign/mgo"
	"github.com/romanornr/cyberchain/insightjson"
)

type MgoBlockDAO struct {
	databaseName string
	session      *mgo.Session
}

// NewBlockDAO creates a new MgoBlockDAO
func NewBlockDAO(database dbName) *MgoBlockDAO {
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

	session.SetMode(mgo.Monotonic, true)

	return &MgoBlockDAO{
		session:      session,
		databaseName: string(database),
	}
}

// DropDatabase drops the entire database.
// Use with caution, supposed to be used only in tests.
func (dao *MgoBlockDAO) DropDatabase() error {
	return dao.session.DB(dao.databaseName).DropDatabase()
}

func (dao *MgoBlockDAO) Get(hash *chainhash.Hash) (*insightjson.BlockResult, error) {
	// reading may be slow, so open extra session here
	session := dao.session.Clone()
	defer session.Close()

	collection := session.DB(dao.databaseName).C(blocks)

	result := &insightjson.BlockResult{}

	err := collection.FindId(hash.String()).One(result)
	if err != nil {
		return nil, err
	}

	err = result.Validate()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (dao *MgoBlockDAO) Create(block *insightjson.BlockResult) error {
	// i guess no need in extra session
	collection := dao.session.DB(dao.databaseName).C(blocks)

	err := collection.Insert(block)
	if err != nil {
		return err
	}

	return nil
}

func (dao *MgoBlockDAO) Delete(hash *chainhash.Hash) error {
	collection := dao.session.DB(dao.databaseName).C(blocks)

	err := collection.RemoveId(hash.String())
	if err != nil {
		return err
	}

	return nil
}
