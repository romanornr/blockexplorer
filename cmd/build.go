// Copyright (c) 2019 Romano, Viacoin developer
//
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/romanornr/cyberchain/blockdata"
	"github.com/romanornr/cyberchain/cmd/rebuilddb"
	"github.com/romanornr/cyberchain/mongodb"
)

var dao = mongodb.MongoDAO{
	"127.0.0.1",
	"viacoin",
}

// This function can be exected with: go run cmd/build.go
// this will build the entire database with blocks, transactions etc
func main() {

	dao.Connect()
	dao.DropDatabase() // delete existing database first

	tip, err := blockdata.GetLatestBlock()
	if err != nil {
		logrus.Fatal("could not get the tip/latest block of the chain")
	}
	rebuilddb.BuildDatabase(tip.Height)
}
