package daos_test

import (
	"github.com/globalsign/mgo"
	. "github.com/romanornr/cyberchain/daos"
	"github.com/romanornr/cyberchain/testdata"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewAddrDAO(t *testing.T) {
	viaAddrDAO := NewAddrDAO(ViaTest)
	viaAddrDAO.DropDatabase() // Ensure everything is clear.
	defer viaAddrDAO.DropDatabase() // Clean up after the test runs.

	// Put address into DB.
	if err := viaAddrDAO.Create(testdata.AddrStdVLavInfo); err != nil {
		t.Fatal(err)
	}

	// Get the address from DB.
	addrID := testdata.AddrStdVLavInfo.Address
	result, err := viaAddrDAO.Get(addrID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, testdata.AddrStdVLavInfo, result)

	// Delete from DB.
	if err = viaAddrDAO.Delete(addrID); err != nil {
		t.Fatal(err)
	}
	// Check deletion.
	if _, err = viaAddrDAO.Get(addrID); err != mgo.ErrNotFound {
		t.Fatal(err)
	}
}
