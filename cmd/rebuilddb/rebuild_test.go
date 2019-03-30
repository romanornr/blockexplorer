// Copyright (c) 2019 Romano, Viacoin developer
//
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package rebuilddb

import (
	"github.com/romanornr/cyberchain/mongodb"
	"log"
	"testing"
)

var dao = mongodb.MongoDAO{
	"127.0.0.1",
	"viacoin",
}

func TestDropDatabase(t *testing.T) {
	dao.Connect()

	log.Println("Dropping old existing database")
	dao.DropDatabase()

	//databases, _ := db.Session.DatabaseNames()
	//
	//for _, databases := range databases {
	//	if databases == "Viacoin" {
	//		fmt.Println("found")
	//		t.Error("Old database still exists. Failed dropping.")
	//	}
	//}
	//log.Println("Success dropped old database")
}

func TestBuildDatabase(t *testing.T) {
	mockTipHeight := int64(11139 + 1)
	BuildDatabase(mockTipHeight)
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

func BenchmarkBuildDatabase(b *testing.B) {
	dao.DropDatabase()
	BuildDatabase(2000)
}
