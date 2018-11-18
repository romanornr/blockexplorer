package subsidy

import (
	"testing"
)

func TestCalcViacoinBlockSubsidy(t *testing.T) {
	valueOut := CalcViacoinBlockSubsidy(5647748, true)
	if valueOut != 0.078125 {
		t.Errorf("Error value out. Expected: 0.078125 Got: %f", valueOut)
	} else {
		t.Logf("Success: blockheight 5647748 got valueOut: %f", valueOut)
	}
}
