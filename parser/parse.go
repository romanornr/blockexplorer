package parser

import (
	"os/user"
	"log"
	"path/filepath"
	"sort"
	"os"
)

// get user home directory
// example: /home/romano
func HomeDirectory() string {
	user, err := user.Current()
	if err != nil {
		log.Fatalf("Error setting home dir to read blocks: %v", err)
	}
	return user.HomeDir
}

// get directory where the blocks are stored on the OS. Assuming coin's default dit
// example: /home/romano/.viacoin/blocks
func BlocksDirectory() string {
	return filepath.Join(HomeDirectory(), ".viacoin", "blocks")
}

// get all the blk.dat files from the Blocks directory
// this will return an array of strings
// example files[0]: /home/romano/.viacoin/blocks/blk00000.dat
func BlockFiles() []string {
	files, _ := filepath.Glob(filepath.Join(BlocksDirectory(), "blk*.dat"))
	sort.Strings(files)
	return files
}

//type Block struct {
//	length            uint32
//	version           uint32
//	hash              [32]byte
//	previousBlockHash [32]byte
//	merkleRoot        [32]byte
//	timestamp         uint32
//	difficulty        uint32
//	nonce             uint32
//	transactionCount  uint64
//}

func BlocksReader() {
	blkFiles := BlockFiles()
	readers := make([]*os.File, len(blkFiles))

	for i, blkFile := range blkFiles {
		reader, err := os.Open(blkFile)
		if err != nil {
			log.Fatalf("Error reading blk file: %v", err)
		}
		readers[i] = reader
	}
}