package daos_test

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/globalsign/mgo"
	. "github.com/romanornr/cyberchain/daos"
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
		t.Fatal(err)
	}

	result, err := viaTxDAO.Get(txID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, testdata.TxBlock5421176Insight, result)

	// Delete from DB.
	if err = viaTxDAO.Delete(txID); err != nil {
		t.Fatal(err)
	}
	// Check deletion.
	if _, err = viaTxDAO.Get(txID); err != mgo.ErrNotFound {
		t.Fatal(err)
	}
}
