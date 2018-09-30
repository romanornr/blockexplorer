package database

import (
	"github.com/romanornr/cyberchain/client"
	"github.com/spf13/viper"
	"os"
	"path"
	"runtime"
	"testing"
)

func TestOpen(t *testing.T) {
	Open()
	client.GetViperConfig()

	_, filename, _, _ := runtime.Caller(0)                    // get full path of this file
	coinsymbol := viper.GetString("coin.symbol")              // example: btc or via
	dbfile := path.Join(path.Dir(filename), coinsymbol+".db") //btc.db or via.db

	// check if file exist
	_, err := os.Stat(dbfile)
	if err != nil {
		t.Error("database file does not exist. Example: via.db")
	}
}

func TestSetupDB(t *testing.T) {
	_, err := SetupDB()
	if err != nil {
		t.Errorf("Error by seting up database buckets: %s", err)
	}
}
