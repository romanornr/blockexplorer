package daos_test

import (
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/globalsign/mgo"
	. "github.com/romanornr/cyberchain/daos"
	"github.com/romanornr/cyberchain/testdata"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewBlockDAO(t *testing.T) {
	viaBlockDAO := NewBlockDAO(ViaTest)
	viaBlockDAO.DropDatabase()       // Ensure everything is clear.
	defer viaBlockDAO.DropDatabase() // Clean up after the test runs.

	// Put test transaction into DB.
	if err := viaBlockDAO.Create(testdata.Block2); err != nil {
		t.Fatal(err)
	}

	// Get chain hash of test transaction for the next lookup.
	blockID, err := chainhash.NewHashFromStr(testdata.Block2.Hash)
	if err != nil {
		t.Fatal(err)
	}

	result, err := viaBlockDAO.Get(blockID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, testdata.Block2, result)

	// Delete from DB.
	if err = viaBlockDAO.Delete(blockID); err != nil {
		t.Fatal(err)
	}
	// Check deletion.
	if _, err = viaBlockDAO.Get(blockID); err != mgo.ErrNotFound {
		t.Fatal(err)
	}
}
