package main

import (
	"fmt"
	"github.com/coreos/bbolt"
	"github.com/romanornr/cyberchain/blockdata"
	"github.com/spf13/viper"
	"log"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"runtime"
	"path"
	"time"
	"encoding/binary"
)

// initalize and read viper configuration
// Create or Open database with the Open() function
// Setup the databse with the SetupDB function
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

// Open or Create a databse in cmd/rebuilddb directory
func Open() error {
	var err error
	_, filename, _, _ := runtime.Caller(0)  // get full path of this file
	coinsymbol := viper.GetString("coin.symbol") // example: btc or via
	dbfile := path.Join(path.Dir(filename), coinsymbol+".db")
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

var Blocks *chainhash.Hash

// add a block
func addBlock(db *bolt.DB, blockheight uint64, block *chainhash.Hash) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Blocks"))

		buf := make([]byte, 8)
		binary.LittleEndian.PutUint64(buf, blockheight) // to byte little endian

		err := b.Put(buf, []byte(block.CloneBytes()))
		return err
	})
}

func viewBlock(blockheight int64) error{
	return db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("Blocks"))
		if bucket == nil{
			return fmt.Errorf("Bucket not found")
		}
		block := bucket.Get([]byte("0"))
		fmt.Println(block)
		return nil
	})
}

func main() {
	blockHeight := int64(0)
	block := blockdata.GetBlockHash(blockHeight)

	fmt.Println("adding block")
	addBlock(db, uint64(blockHeight), block)
	viewBlock(blockHeight)

}
