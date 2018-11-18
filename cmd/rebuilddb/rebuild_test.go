package rebuilddb

import (
	"testing"
	"github.com/romanornr/cyberchain/mongodb"
	"log"
	"fmt"
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

//func BenchmarkBuildDatabase(b *testing.B) {
//	BuildDatabase()
//}
