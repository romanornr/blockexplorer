package daos

import (
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/globalsign/mgo"
	"github.com/romanornr/cyberchain/insightjson"
)

type MgoTxDAO struct {
	databaseName string
	session      *mgo.Session
}

// NewBlockDAO creates a new MgoBlockDAO
func NewTxDAO(database dbName) *MgoTxDAO {
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

	return &MgoTxDAO{
		session:      session,
		databaseName: string(database),
	}
}

// DropDatabase drops the entire database.
// Use with caution, supposed to be used only in tests.
func (dao *MgoTxDAO) DropDatabase() error {
	return dao.session.DB(dao.databaseName).DropDatabase()
}

func (dao *MgoTxDAO) Get(hash *chainhash.Hash) (*insightjson.Tx, error) {
	// reading may be slow, so open extra session here
	session := dao.session.Clone()
	defer session.Close()

	collection := session.DB(dao.databaseName).C(txs)

	result := &insightjson.Tx{}

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

func (dao *MgoTxDAO) Create(tx *insightjson.Tx) error {
	// i guess no need in extra session
	collection := dao.session.DB(dao.databaseName).C(txs)

	index := mgo.Index{
		Key:    []string{"txid"},
		Unique: true,
	}

	err := collection.EnsureIndex(index)
	if err != nil {
		return err
	}

	err = collection.Insert(tx)
	if err != nil {
		return err
	}

	return nil
}

func (dao *MgoTxDAO) Delete(hash *chainhash.Hash) error {
	collection := dao.session.DB(dao.databaseName).C(txs)

	err := collection.RemoveId(hash.String())
	if err != nil {
		return err
	}

	return nil
}
