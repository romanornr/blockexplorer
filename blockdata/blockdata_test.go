package blockdata

import (
	"fmt"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"log"
	"testing"
)

func TestGetBlockCount(t *testing.T) {
	result := GetBlockCount()
	if result < 1 {
		t.Errorf("Expected the blockcount to be higher than 1")
	}
	log.Printf("Success: Current blockcount is %d\n", result)
}

func TestGetBlockHash(t *testing.T) {
	result, _ := GetBlockHash(1)
	expected := "5ca83af67146e286610e118cc8f8e6a183c319fbb4a8fdb9e99daa2b8a29b3e3"
	if result.String() != expected {
		t.Errorf("Expected ")
	}
	log.Printf("Success: blockheight 1 is equal to %s", result.String())
}

func TestGetBlock(t *testing.T) {
	blockhash, _ := GetBlockHash(1)
	result, _ := GetBlock(blockhash)
	fmt.Println(result)
}

func TestGetRawTransactionVerbose(t *testing.T) {
	hash, _ := chainhash.NewHashFromStr("36a8207c5a3c46974ef777f5904c62409137eda73103ca324deb82b894785e9d")
	transaction, _ := GetRawTransactionVerbose(hash)
	fmt.Println(transaction)
}
