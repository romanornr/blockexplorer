package database

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"log"
	"path"
	"runtime"
	"time"

	"errors"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/coreos/bbolt"
	"github.com/romanornr/cyberchain/blockdata"
	"github.com/romanornr/cyberchain/client"
	"github.com/spf13/viper"
)

var db *bolt.DB

// FIXME: Create more static errors for the project.
var (
	errBucketNotFound = errors.New("bucket is not found")
)

// initialize and read viper configuration
// create or Open database with the Open() function
// setup the database with the SetupDB function
func GetDatabaseInstance() *bolt.DB {
	if db != nil {
		return db
	}
	client.GetViperConfig()
	Open()
	SetupDB()
	// instance := new(*bolt.DB)
	fmt.Println("new instance")

	return db
}

// open or Create a database in cmd/rebuilddb directory
func Open() error {
	var err error
	_, filename, _, _ := runtime.Caller(0)                    // get full path of this file
	coinsymbol := viper.GetString("coin.symbol")              // example: btc or via
	dbfile := path.Join(path.Dir(filename), coinsymbol+".db") // btc.db or via.db
	config := &bolt.Options{Timeout: 10 * time.Second}
	db, err = bolt.Open(dbfile, 0600, config)
	if err != nil {
		log.Fatal(err)
	}

	return err
}

// setup database with a bucket called Blocks
func SetupDB() (*bolt.DB, error) {
	var err error

	err = db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte("Blocks"))
		if err != nil {
			return fmt.Errorf("could not create blocks bucket: %v", err)
		}
		return nil
	})

	err = db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte("BlockHeight"))
		if err != nil {
			return fmt.Errorf("could not create blockheight bucket: %v", err)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("could not setup buckets, %v", err)
	}

	fmt.Println("DB Setup Done")

	return db, nil
}

// add a block to the database and use the CloneBytes() function to put the blocks to byte.
// the blockhash string is the key. The value is all the data in the block
func AddBlock(db *bolt.DB, blockHashString string, block *btcjson.GetBlockVerboseResult) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Blocks"))

		var result bytes.Buffer
		encoder := gob.NewEncoder(&result)
		encoder.Encode(block)

		// check if the previous blockheight is not higher than the current blockheight.
		prevBlockHash, _ := chainhash.NewHashFromStr(block.PreviousHash)
		prevBlockHeader := blockdata.GetBlockHeader(prevBlockHash)

		if int32(block.Height) < prevBlockHeader.Height {
			log.Panic("Error: Previous blockheight is higher than the current blockheight. Something went wrong.")
		}

		return b.Put([]byte(blockHashString), result.Bytes())
	})
}

// link in the boltdb database the blockheight with the right blockhash.
// this way the blockheight can be used to find the right blockhash.
func AddIndexBlockHeightWithBlockHash(db *bolt.DB, blockHash string, blockHeight int64) error {
	return db.Update(func(tx *bolt.Tx) error {
		b, _ := tx.CreateBucketIfNotExists([]byte("Blockheight"))

		// blockheight to byte
		bs := make([]byte, 8)
		binary.BigEndian.PutUint64(bs, uint64(blockHeight))

		return b.Put([]byte(bs), []byte(blockHash))
	})
}

// link in botldb database the transaction with the right blockhash
func AddTransaction(db *bolt.DB, TransactionHash []string) error {
	return db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("Transactions"))
		if err != nil {
			log.Fatalf("Error: Could not save transaction together with the blockhash in the Transaction bucket: %v", err)
		}

		for _, element := range TransactionHash {
			txHash, _ := chainhash.NewHashFromStr(element)
			rawTransaction := blockdata.GetRawTransactionVerbose(txHash)
			// rawTransaction := blockdata.GetRawTransaction(txHash)

			var result bytes.Buffer
			encoder := gob.NewEncoder(&result)
			encoder.Encode(rawTransaction.Hex)
			b.Put(txHash.CloneBytes(), []byte(rawTransaction.Hex))
		}
		return nil
	})
}

// view the block by giving the blockhash string
// Is that a good idea to return just a slice?
// Maybe it should be a Block struct, so the func will return
// (*btcjson.GetBlockVerboseResult, err) or
// (*btcjson.GetBlockVerboseResult, bool), where bool is existence of the block.
func ViewBlock(blockHashString string) []byte {
	var block []byte
	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("Blocks"))
		if bucket == nil {
			return errBucketNotFound
		}

		block = bucket.Get([]byte(blockHashString))
		return nil
	})
	return block
}

func FetchTransactionHashByBlockhash(blockHashString string) []byte {
	var block []byte
	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("Transactions"))
		if bucket == nil {
			return errBucketNotFound
		}

		block = bucket.Get([]byte(blockHashString))
		return nil
	})
	return block
}

func FetchBlockHashByBlockHeight(blockheight int64) []byte {
	var hash []byte
	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("Blockheight"))
		if bucket == nil {
			return errBucketNotFound
		}

		bs := make([]byte, 8)
		binary.BigEndian.PutUint64(bs, uint64(blockheight))

		hash = bucket.Get([]byte(bs))
		return nil
	})
	fmt.Println(hash)
	return hash
}
