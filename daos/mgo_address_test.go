package daos_test

import (
	. "github.com/romanornr/cyberchain/daos"
	"testing"
)

func TestNewAddrDAO(t *testing.T) {
	viaAddrDAO := NewAddrDAO(ViaTest)
	viaAddrDAO.DropDatabase() // Ensure everything is clear.
	defer viaAddrDAO.DropDatabase() // Clean up after the test runs.

	// Put address into DB.
}
