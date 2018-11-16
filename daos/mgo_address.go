package daos

import (
	"github.com/globalsign/mgo"
	"github.com/romanornr/cyberchain/insightjson"
)

type MgoAddrDAO struct {
	databaseName string
	session      *mgo.Session
}

// NewAddrDAO creates a new MgoBlockDAO
func NewAddrDAO(database dbName) *MgoAddrDAO {
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

	return &MgoAddrDAO{
		session:      session,
		databaseName: string(database),
	}
}

// DropDatabase drops the entire database.
// Use with caution, supposed to be used only in tests.
func (dao *MgoAddrDAO) DropDatabase() error {
	return dao.session.DB(dao.databaseName).DropDatabase()
}

func (dao *MgoAddrDAO) Get(addrID string) (*insightjson.AddressInfo, error) {
	// reading may be slow, so open extra session here
	session := dao.session.Clone()
	defer session.Close()

	collection := session.DB(dao.databaseName).C(addr)

	result := &insightjson.AddressInfo{}

	err := collection.FindId(addrID).One(result)
	if err != nil {
		return nil, err
	}

	err = result.Validate()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (dao *MgoAddrDAO) Create(address *insightjson.AddressInfo) error {
	// i guess no need in extra session
	collection := dao.session.DB(dao.databaseName).C(addr)

	err := collection.Insert(address)
	if err != nil {
		return err
	}

	return nil
}

func (dao *MgoAddrDAO) Delete(addrID string) error {
	collection := dao.session.DB(dao.databaseName).C(addr)

	err := collection.RemoveId(addrID)
	if err != nil {
		return err
	}

	return nil
}
