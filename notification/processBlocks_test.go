package notification

import (
	"testing"
)

// test if pools.json can be read
// and contains the correct content
func TestParseJson(t *testing.T) {
	ParseJson()
	result := pools[0].PoolName
	expected := "50BTC"
	if result != expected {
		t.Errorf("Got %s, expected %s", result, expected)
	}
}
