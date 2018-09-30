package client

import (
	"fmt"
	"testing"
)

func TestGetInstance(t *testing.T) {
	instance1 := GetInstance()
	if instance1 == nil {
		t.Error("Expetected pointer to Singleton after calling GetInstance(), not nill")
	}
	fmt.Println(instance1.GetBlockCount())
}
