package main

import (
	"testing"
	"github.com/romanornr/cyberchain/database"
)

func Testmain(t *testing.T) {
	//now := time.Now().UTC()
	database.BuildDatabaseBlocks()
}