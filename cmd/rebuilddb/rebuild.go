package main

import (
	"time"
	"fmt"
	"github.com/romanornr/cyberchain/database"
)

func main() {
	now := time.Now().UTC()
	database.BuildDatabaseBlocks()
	fmt.Println("Time elapsed: ", time.Since(now))
}
