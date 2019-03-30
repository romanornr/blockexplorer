// Copyright (c) 2019 Romano, Viacoin developer
//
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package mongodb

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcutil"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/romanornr/cyberchain/insightjson"
)

type MongoDAO struct {
	Server   string
	Database string
}

var db *mgo.Database
var dialInfo = viaDialInfo

func (dao *MongoDAO) Connect() {
	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		session.Close()
		log.Panicf("failed to open mongodb session: %s\n", err)
	}

	session.SetMode(mgo.Monotonic, true)
	db = session.DB(viaDialInfo.Database)
}


func (dao *MongoDAO) GetBlock(hash *chainhash.Hash) (*insightjson.BlockResult, error) {
	dao.Connect()
	defer db.Session.Close()

	result := new(insightjson.BlockResult)

	err := db.C(BLOCKS).Find(bson.M{"hash": hash.String()}).One(&result)
	if err != nil {
		return result, err
	}

	return result, err
}


func (dao *MongoDAO) AddBlock(Block *insightjson.BlockResult) error {

	dao.Connect()
	defer db.Session.Close()

	index := mgo.Index{
		Key:    []string{"hash"},
		Unique: true,
	}

	err := db.C(BLOCKS).EnsureIndex(index)
	if err != nil {
		log.Panicf("%s\n", err)
	}

	err = db.C(BLOCKS).Insert(Block)

	if err != nil {
		log.Warningf("Block with hash %s did not get inserted\n", Block.Hash)
		return fmt.Errorf("Block with hash %s did not get inserted\n", Block.Hash)
	}

	return err
}

// get the latest block and return it in insightjson format
func (dao *MongoDAO) GetLastBlock() (insightjson.BlockResult, error) {
	dao.Connect()
	defer db.Session.Close()

	result := insightjson.BlockResult{}

	dbSize, err := db.C(BLOCKS).Count()
	if err != nil {
		return result, err
	}

	err = db.C(BLOCKS).Find(nil).Skip(dbSize - 1).One(&result)
	if err != nil {
		return result, err
	}

	return result, err
}

// add transaction to the database. Make sure the txid is unique.
func (dao *MongoDAO) AddTransaction(transaction *insightjson.Tx) error {
	dao.Connect()
	defer db.Session.Close()

	index := mgo.Index{
		Key:    []string{"txid"},
		Unique: true,
	}

	err := db.C(TRANSACTIONS).EnsureIndex(index)
	if err != nil {
		panic(err)
	}

	err = db.C(TRANSACTIONS).Insert(transaction)

	if err != nil {
		panic(err)
	}

	return err
}

// Find transactions by txid in the database
func (dao *MongoDAO) GetTransaction(txid chainhash.Hash) (insightjson.Tx, error) {
	dao.Connect()
	defer db.Session.Close()

	result := insightjson.Tx{}

	err := db.C(TRANSACTIONS).Find(bson.M{"txid": txid.String()}).One(&result)
	if err != nil {
		return result, err
	}

	return result, err
}


// delete the database. Only use for testing
func (dao *MongoDAO) DropDatabase() error {
	dao.Connect()
	//defer db.Session.Close()

	err := db.Session.DB(db.Name).DropDatabase()
	//err := session.DB(Database).DropDatabase()
	if err != nil {
		panic(err)
	}
	return err
}

// add addressInfo to the database
func (dao *MongoDAO) AddAddressInfo(AddressInfo *insightjson.AddressInfo) error {
	dao.Connect()
	defer db.Session.Close()


	err := db.C(ADDRESSINFO).Insert(AddressInfo)
	if err != nil {
		log.Printf("Error not able to add AddressInfo to database collection AddressInfo: %s", err)
	}

	return nil
}

func (dao *MongoDAO) GetAddressInfo(address string) (insightjson.AddressInfo, error) {
	dao.Connect()
	defer db.Session.Close()

	result := insightjson.AddressInfo{}

	err := db.C(ADDRESSINFO).Find(bson.M{"addrStr": address}).One(&result)
	if err != nil {
		return result, err
	}

	return result, err
}

func (dao *MongoDAO) UpdateAddressInfoSent(AddressInfo *insightjson.AddressInfo, sentSat int64, confirmed bool, txid string) error {

	dao.Connect()
	defer db.Session.Close()

	colQuerier := bson.M{"addrStr": AddressInfo.Address}

	AddressInfo.TransactionsID = append(AddressInfo.TransactionsID, txid) //TODO: change order

	if !confirmed {
		AddressInfo.UnconfirmedTxAppearances += 1
		AddressInfo.UnconfirmedBalance -= btcutil.Amount(sentSat).ToBTC()
		AddressInfo.UnconfirmedBalanceSat = sentSat
		change := bson.M{"$set": bson.M{"unconfirmedTxAppearances": AddressInfo.UnconfirmedTxAppearances, "unconfirmedBalance": AddressInfo.UnconfirmedBalance, "unconfirmedBalanceSat": AddressInfo.UnconfirmedBalanceSat, "transactions": AddressInfo.TransactionsID}}
		err := db.C(ADDRESSINFO).Update(colQuerier, change)
		if err != nil {
			log.Printf("Failed to update AddressInfo for adress %s: %s", AddressInfo.Address, err)
		}
		return err
	}

	AddressInfo.TxAppearances += 1
	AddressInfo.TotalSentSat += sentSat
	AddressInfo.TotalSent += btcutil.Amount(sentSat).ToBTC()
	AddressInfo.BalanceSat -= sentSat
	AddressInfo.Balance -= btcutil.Amount(sentSat).ToBTC()

	change := bson.M{"$set": bson.M{"totalSentSat": AddressInfo.TotalSentSat, "totalSent": AddressInfo.TotalSent, "txAppearances": AddressInfo.TxAppearances, "balanceSat": AddressInfo.BalanceSat, "balance": AddressInfo.Balance, "transactions": AddressInfo.TransactionsID}}
	err := db.C(ADDRESSINFO).Update(colQuerier, change)
	if err != nil {
		log.Printf("Failed to update AddressInfo for adress %s: %s", AddressInfo.Address, err)
	}
	return err //TODO What if it's still unconfirmed. Unconfirmed Balance & Unconfirmed TotalSent & Unconfirmed tx Appearances
}

// Update addressInfo by searching with the address string
func (dao *MongoDAO) UpdateAddressInfoReceived(AddressInfo *insightjson.AddressInfo, receivedSat int64, confirmed bool, txid string) error {

	dao.Connect()
	defer db.Session.Close()

	colQuerier := bson.M{"addrStr": AddressInfo.Address}

	AddressInfo.TransactionsID = append(AddressInfo.TransactionsID, txid) //TODO change order

	if !confirmed {
		AddressInfo.UnconfirmedTxAppearances += 1
		AddressInfo.UnconfirmedBalance += btcutil.Amount(receivedSat).ToBTC()
		AddressInfo.UnconfirmedBalanceSat += receivedSat
		change := bson.M{"$set": bson.M{"unconfirmedTxAppearances": AddressInfo.UnconfirmedTxAppearances, "unconfirmedBalance": AddressInfo.UnconfirmedBalance, "unconfirmedBalanceSat": AddressInfo.UnconfirmedBalanceSat, "transactions": AddressInfo.TransactionsID}}
		err := db.C(ADDRESSINFO).Update(colQuerier, change)
		if err != nil {
			log.Printf("Failed to update AddressInfo for adress %s: %s", AddressInfo.Address, err)
		}
		return err
	}

	AddressInfo.TxAppearances += 1
	AddressInfo.TotalReceivedSat += receivedSat
	AddressInfo.TotalReceived += btcutil.Amount(receivedSat).ToBTC()
	AddressInfo.BalanceSat += receivedSat
	AddressInfo.Balance += btcutil.Amount(receivedSat).ToBTC()

	change := bson.M{"$set": bson.M{"totalReceivedSat": AddressInfo.TotalReceivedSat, "totalReceived": AddressInfo.TotalReceived, "txAppearances": AddressInfo.TxAppearances, "balanceSat": AddressInfo.BalanceSat, "balance": AddressInfo.Balance, "transactions": AddressInfo.TransactionsID}}
	err := db.C(ADDRESSINFO).Update(colQuerier, change)
	if err != nil {
		log.Printf("Failed to update AddressInfo for adress %s: %s", AddressInfo.Address, err)
	}
	return err //TODO What if it's still unconfirmed. Unconfirmed Balance & Unconfirmed TotalSent & Unconfirmed tx Appearances
}

// Update Transaction by txid
func (dao *MongoDAO) UpdateTransaction(tx *insightjson.Tx) error {
	dao.Connect()
	defer db.Session.Close()

	selector := bson.M{"txid": tx.Txid}

	_, err := db.C(TRANSACTIONS).Upsert(selector, tx)
	if err != nil {
		log.Printf("Error updating spentDetails in mongodb: %s", err)
	}

	return err
}

// Get the unspent outs of an address
// by getting the addressInfo of an address and aggregate over all transaction id's from the past
// get the tx info from the database and check for the vouts
// example: https://explorer.viacoin.org/api/addr/VmkyKgGBWDpcnFCtw8rvcYqg8xr7U4Ubzx/utxo
func (dao *MongoDAO) GetAddressUTXO(address string) []insightjson.UnpsentOutput {

	dao.Connect()
	defer db.Session.Close()

	utxo := []insightjson.UnpsentOutput{}

	addressInfo, err := dao.GetAddressInfo(address)
	if err != nil {
		return utxo
	}

	txHistory := addressInfo.TransactionsID

	for _, h := range txHistory {
		hash, err := chainhash.NewHashFromStr(h)
		if err != nil {
			fmt.Errorf("Not able to make tx hash into chainhash: %s", err)
			return utxo
		}

		tx, err := dao.GetTransaction(*hash)
		if err != nil {
			fmt.Errorf("failed to get transaction from mongodb: %s", err)
		}

		for idx, vout := range tx.Vouts {
			//	if vout.SpentIndex != nil {  // TODO unsure yet. Only show non-empty utxo's or non-empty and empty. For now showing both
			satoshi := btcutil.Amount(vout.Value)

			output := insightjson.UnpsentOutput{
				Address:      address,
				Txid:         tx.Txid,
				Vout:         idx,
				ScriptPubKey: vout.ScriptPubKey.Hex,
				Amount:       vout.Value,
				Satoshis:     int64(satoshi),
				Height:       tx.Blockheight,
			}
			utxo = append(utxo, output)
		}
	}
	//}
	return utxo
}
