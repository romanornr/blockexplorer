package subsidy

import (
	"testing"
	"github.com/viacoin/viad/chaincfg"
)

func TestCalcViacoinBlockSubsidy(t *testing.T) {
	valueOut := CalcViacoinBlockSubsidy(5647748, &chaincfg.MainNetParams)
	if valueOut != 7812500 {
		t.Errorf("Error value out. Expected: 7812500 expected: %d", valueOut)
	}else {
		t.Logf("Success: blockheight 5647748 got valueOut %d", valueOut)
	}
}
