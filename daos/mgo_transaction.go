package daos

// type MgoTxDAO struct {
// 	databaseName string
// 	session      *mgo.Session
// }
//
// // NewBlockDAO creates a new MgoBlockDAO
// func NewTxDAO(database dbName) *MgoTxDAO {
// 	session, err := mgo.DialWithInfo(
// 		&mgo.DialInfo{
// 			Addrs:    dbHosts,
// 			Database: string(database),
// 			Timeout:  timeout10,
// 		},
// 	)
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	session.SetMode(mgo.Monotonic, true)
//
// 	return &MgoTxDAO{
// 		session:      session,
// 		databaseName: string(database),
// 	}
// }
//
// func (dao *MgoTxDAO) Get(hash chainhash.Hash) (*insightjson.BlockResult, error) {
// 	// reading may be slow, so open extra session here
// 	session := dao.session.Clone()
// 	defer session.Close()
//
// 	collection := session.DB(dao.databaseName).C(blocks)
//
// 	result := &insightjson.BlockResult{}
//
// 	err := collection.FindId(hash.String()).One(result)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	err = result.Validate()
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	return result, nil
// }
//
// func (dao *MgoTxDAO) Create(block *insightjson.BlockResult) error {
// 	// i guess no need in extra session
// 	collection := dao.session.DB(dao.databaseName).C(blocks)
//
// 	index := mgo.Index{
// 		Key:    []string{"hash"},
// 		Unique: true,
// 	}
//
// 	err := collection.EnsureIndex(index)
// 	if err != nil {
// 		return err
// 	}
//
// 	err = collection.Insert(block)
// 	if err != nil {
// 		return err
// 	}
//
// 	return nil
// }
//
// func (dao *MgoTxDAO) Delete(hash chainhash.Hash) error {
// 	// maybe need to use this
// 	// _, err := dao.Get(hash)
// 	// if err != nil {
// 	// 	// no block exists
// 	// }
//
// 	// i guess no need in extra session
// 	// however it performs reading, idk
// 	collection := dao.session.DB(dao.databaseName).C(blocks)
// 	err := collection.RemoveId(hash)
// 	if err != nil {
// 		return err
// 	}
//
// 	return nil
// }
//
