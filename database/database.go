package database

import (
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcd/wire"
	"github.com/coreos/bbolt"
	"github.com/romanornr/cyberchain/blockdata"
	"github.com/spf13/viper"
	"log"
	"path"
	"runtime"
	"time"
)

// initalize and read viper configuration
// create or Open database with the Open() function
// setup the database with the SetupDB function
func init() {

	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")
	viper.SetConfigName("app")

	err := viper.ReadInConfig()

	if err != nil {
		log.Fatal("No configuration file loaded ! Please check the config folder")
	}

	fmt.Printf("Reading configuration from %s\n", viper.ConfigFileUsed())
	Open()
	setupDB()
}

var db *bolt.DB
var open bool

// open or Create a databse in cmd/rebuilddb directory
func Open() error {
	var err error
	_, filename, _, _ := runtime.Caller(0)       // get full path of this file
	coinsymbol := viper.GetString("coin.symbol") // example: btc or via
	dbfile := path.Join(path.Dir(filename), coinsymbol+".db") //btc.db or via.db
	config := &bolt.Options{Timeout: 1 * time.Second}
	db, err = bolt.Open(dbfile, 0600, config)
	if err != nil {
		log.Fatal(err)
	}
	open = true
	return err
}

// setup database with a bucket called Blocks
func setupDB() (*bolt.DB, error) {
	var err error

	err = db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte("Blocks"))
		if err != nil {
			return fmt.Errorf("could not create blocks bucket: %v", err)
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
func AddBlock(db *bolt.DB, blockHashString string, block *wire.MsgBlock) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Blocks"))
		encoded, err := json.Marshal(block)
		if err != nil {
			log.Fatalf("could not add block to database: %v", err)
		}

		return b.Put([]byte(blockHashString), encoded)
	})
}

// view the block by giving the blockhash string
func ViewBlock(blockHashString string) []byte {
	var block []byte
	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("Blocks"))
		if bucket == nil {
			return fmt.Errorf("bucket not found")
		}

		block = bucket.Get([]byte(blockHashString))
		return nil
	})
	return block
}

func BuildDatabaseBlocks()  {
	for i := int64(0); i < 100; i++ {
		height := blockdata.GetBlockHash(i)
		block := blockdata.GetBlock(height)
		blockHashString := block.Header.BlockHash().String()
		AddBlock(db, blockHashString, block)
	}
}

