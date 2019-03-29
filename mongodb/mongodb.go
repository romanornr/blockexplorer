// Copyright (c) 2019 Romano, Viacoin developer
//
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package mongodb

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"time"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcutil"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/romanornr/cyberchain/insightjson"
)

const (
	MongoDBHosts = "localhost"
	Database     = "viacoin"
)

var session *mgo.Session

var mongoDBDialInfo = &mgo.DialInfo{
	Addrs:    []string{MongoDBHosts},
	Timeout:  60 * time.Second,
	Database: Database,
}

func GetSession() *mgo.Session {
	if session == nil {
		var err error
		session, err = mgo.DialWithInfo(mongoDBDialInfo)
		if err != nil {
			session.Close()
			panic(err)
		}
	}
	return session
}

// delete the database. Only use for testing
func DropDatabase() error {
	err := session.DB(Database).DropDatabase()
	if err != nil {
		panic(err)
	}
	return err
}

// add Blocks to the database. Collection name: Blocks
func AddBlock(Block *insightjson.BlockResult) error {
	GetSession()
	//defer session.Close()

	collection := session.DB(Database).C("Blocks")

	index := mgo.Index{
		Key:    []string{"hash"},
		Unique: true,
	}

	err := collection.EnsureIndex(index)
	if err != nil {
		panic(err)
	}

	err = collection.Insert(Block)

	if err != nil {
		log.Warningf("Block with hash %s did not get inserted", Block.Hash)
		return fmt.Errorf("Block with hash %s did not get inserted", Block.Hash)
	}

	return err
}

// get block by hash
func GetBlock(hash chainhash.Hash) (insightjson.BlockResult, error) {
	GetSession()
	collection := session.DB(Database).C("Blocks")

	result := insightjson.BlockResult{}

	err := collection.Find(bson.M{"hash": hash.String()}).One(&result)
	if err != nil {
		return result, err
	}

	return result, err
}

// get block by blockheight
func FetchBlockHashByBlockHeight(blockheight int64) insightjson.BlockResult {
	GetSession()
	collection := session.DB(Database).C("Blocks")
	result := insightjson.BlockResult{}

	err := collection.Find(bson.M{"height": blockheight}).One(&result)
	if err != nil {
		panic(err)
	}

	return result
}

// get the latest block and return it in insightjson format
func GetLastBlock() (insightjson.BlockResult, error) {
	GetSession()
	collection := session.DB(Database).C("Blocks")
	result := insightjson.BlockResult{}

	dbSize, err := collection.Count()
	if err != nil {
		return result, err
	}

	err = collection.Find(nil).Skip(dbSize - 1).One(&result)
	if err != nil {
		return result, err
	}

	return result, err
}

func AddTransaction(transaction *insightjson.Tx) error {
	GetSession()
	collection := session.DB(Database).C("Transactions")

	index := mgo.Index{
		Key:    []string{"txid"},
		Unique: true,
	}

	err := collection.EnsureIndex(index)
	if err != nil {
		panic(err)
	}

	err = collection.Insert(transaction)

	if err != nil {
		panic(err)
	}

	return err
}

func GetTransaction(txid chainhash.Hash) (insightjson.Tx, error) {
	GetSession()
	collection := session.DB(Database).C("Transactions")

	result := insightjson.Tx{}

	err := collection.Find(bson.M{"txid": txid.String()}).One(&result)
	if err != nil {
		return result, err
	}

	return result, err
}

func AddAddressInfo(AddressInfo *insightjson.AddressInfo) error {
	GetSession()
	collection := session.DB(Database).C("AddressInfo")

	//if AddressInfo.Address != "" { //check if address is not empty
	err := collection.Insert(AddressInfo)
	if err != nil {
		log.Printf("Error not able to add AddressInfo to database collection AddressInfo: %s", err)
		//}
	}

	return nil
}

func GetAddressInfo(address string) (insightjson.AddressInfo, error) {
	GetSession()
	collection := session.DB(Database).C("AddressInfo")

	result := insightjson.AddressInfo{}

	err := collection.Find(bson.M{"addrStr": address}).One(&result)
	if err != nil {
		return result, err
	}

	return result, err
}

func UpdateAddressInfoSent(AddressInfo *insightjson.AddressInfo, sentSat int64, confirmed bool, txid string) error {

	GetSession()
	collection := session.DB(Database).C("AddressInfo")
	colQuerier := bson.M{"addrStr": AddressInfo.Address}

	AddressInfo.TransactionsID = append(AddressInfo.TransactionsID, txid) //TODO: change order

	if !confirmed {
		AddressInfo.UnconfirmedTxAppearances += 1
		AddressInfo.UnconfirmedBalance -= btcutil.Amount(sentSat).ToBTC()
		AddressInfo.UnconfirmedBalanceSat = sentSat
		change := bson.M{"$set": bson.M{"unconfirmedTxAppearances": AddressInfo.UnconfirmedTxAppearances, "unconfirmedBalance": AddressInfo.UnconfirmedBalance, "unconfirmedBalanceSat": AddressInfo.UnconfirmedBalanceSat, "transactions": AddressInfo.TransactionsID}}
		err := collection.Update(colQuerier, change)
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
	err := collection.Update(colQuerier, change)
	if err != nil {
		log.Printf("Failed to update AddressInfo for adress %s: %s", AddressInfo.Address, err)
	}
	return err //TODO What if it's still unconfirmed. Unconfirmed Balance & Unconfirmed TotalSent & Unconfirmed tx Appearances
}

func UpdateAddressInfoReceived(AddressInfo *insightjson.AddressInfo, receivedSat int64, confirmed bool, txid string) error {

	GetSession()
	collection := session.DB(Database).C("AddressInfo")
	colQuerier := bson.M{"addrStr": AddressInfo.Address}

	AddressInfo.TransactionsID = append(AddressInfo.TransactionsID, txid) //TODO change order

	if !confirmed {
		AddressInfo.UnconfirmedTxAppearances += 1
		AddressInfo.UnconfirmedBalance += btcutil.Amount(receivedSat).ToBTC()
		AddressInfo.UnconfirmedBalanceSat += receivedSat
		change := bson.M{"$set": bson.M{"unconfirmedTxAppearances": AddressInfo.UnconfirmedTxAppearances, "unconfirmedBalance": AddressInfo.UnconfirmedBalance, "unconfirmedBalanceSat": AddressInfo.UnconfirmedBalanceSat, "transactions": AddressInfo.TransactionsID}}
		err := collection.Update(colQuerier, change)
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
	err := collection.Update(colQuerier, change)
	if err != nil {
		log.Printf("Failed to update AddressInfo for adress %s: %s", AddressInfo.Address, err)
	}
	return err //TODO What if it's still unconfirmed. Unconfirmed Balance & Unconfirmed TotalSent & Unconfirmed tx Appearances
}

func UpdateTransaction(tx *insightjson.Tx) error {
	GetSession()
	collection := session.DB(Database).C("Transactions")
	selector := bson.M{"txid": tx.Txid}

	_, err := collection.Upsert(selector, tx)
	if err != nil {
		log.Printf("Error updating spentDetails in mongodb: %s", err)
	}

	return err
}

// Get the unspent outs of an address
// by getting the addressInfo of an address and aggregate over all transaction id's from the past
// get the tx info from the database and check for the vouts
// example: https://explorer.viacoin.org/api/addr/VmkyKgGBWDpcnFCtw8rvcYqg8xr7U4Ubzx/utxo
func GetAddressUTXO(address string) []insightjson.UnpsentOutput {

	utxo := []insightjson.UnpsentOutput{}

	addressInfo, err := GetAddressInfo(address)
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

		tx, err := GetTransaction(*hash)
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
