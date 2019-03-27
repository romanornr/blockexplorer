// Copyright (c) 2019 Romano, Viacoin developer
//
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

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
