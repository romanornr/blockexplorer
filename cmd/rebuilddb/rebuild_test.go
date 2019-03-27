// Copyright (c) 2019 Romano, Viacoin developer
//
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package rebuilddb

import (
	"fmt"
	"github.com/romanornr/cyberchain/mongodb"
	"log"
	"testing"
)

func TestDropDatabase(t *testing.T) {
	mongodb.DropDatabase()
	session := mongodb.GetSession()

	log.Println("Dropping old existing database")
	mongodb.DropDatabase()

	databases, _ := session.DatabaseNames()

	for _, databases := range databases {
		if databases == "Viacoin" {
			fmt.Println("found")
			t.Error("Old database still exists. Failed dropping.")
		}
	}
	log.Println("Success dropped old database")
}

func TestBuildDatabase(t *testing.T) {
	BuildDatabase()
}

//func TestAddrIndex(t *testing.T) {
//
//	address := "Vn5maEzzZNPQ85rKFAgACRW98oiDtmMumG"
//	exepectedBalance := 2360.19356999
//	addressInfo, _ := mongodb.GetAddressInfo(address)
//
//	if addressInfo.Balance != exepectedBalance {
//		t.Errorf("Wrong Balance: expected %f, actual %f", exepectedBalance, addressInfo.Balance)
//	}
//}

//func BenchmarkBuildDatabase(b *testing.B) {
//	mongodb.DropDatabase()
//	BuildDatabase()
//}
