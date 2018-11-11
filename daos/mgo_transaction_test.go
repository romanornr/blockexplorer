package daos_test

import (
	"reflect"
	"testing"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/globalsign/mgo"
	. "github.com/romanornr/cyberchain/daos"
	"github.com/romanornr/cyberchain/insightjson"
	"github.com/romanornr/cyberchain/testdata"
)

func TestMgoTxDAO(t *testing.T) {
	viaTxDAO := NewTxDAO(ViaTest)
	viaTxDAO.DropDatabase()       // Ensure everything is clear.
	defer viaTxDAO.DropDatabase() // Clean up after the test runs.

	// Put test transaction into DB.
	if err := viaTxDAO.Create(testdata.TxBlock5421176Insight); err != nil {
		t.Fatal(err)
	}

	// Get chain hash of test transaction for the next lookup.
	txID, err := chainhash.NewHashFromStr(testdata.TxBlock5421176Insight.Txid)
	if err != nil {
		t.Fatal(txID)
	}

	result, err := viaTxDAO.Get(txID)
	if err != nil {
		t.Fatal(err)
	}
	assertTransactionsEqual(t, testdata.TxBlock5421176Insight, result)

	// Test deletion.
	if err = viaTxDAO.Delete(txID); err != nil {
		t.Fatal(err)
	}
	// Check deletion.
	if _, err = viaTxDAO.Get(txID); err != mgo.ErrNotFound {
		t.Fatal(err)
	}
}

func assertTransactionsEqual(t *testing.T, expected, got *insightjson.Tx) {
	if !reflect.DeepEqual(expected, got) {
		t.Fatalf("expected %v but got %v instead", expected, got)
	}
}
