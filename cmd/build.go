// Copyright (c) 2019 Romano, Viacoin developer
//
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"github.com/romanornr/cyberchain/blockdata"
	"github.com/romanornr/cyberchain/cmd/rebuilddb"
	"github.com/romanornr/cyberchain/mongodb"
	"log"
)

// This function can be exected with: go run cmd/build.go
// this will build the entire database with blocks, transactions etc
func main() {

	mongodb.DropDatabase() // delete existing database first

	tip, err := blockdata.GetLatestBlock()
	if err != nil {
		log.Fatalf("could not get the tip/latest block of the chain")
	}
	rebuilddb.BuildDatabase(tip.Height)
}
