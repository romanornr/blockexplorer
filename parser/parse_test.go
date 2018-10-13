package parser

import (
	"runtime"
	"strings"
	"testing"
	"fmt"
)

func TestHomeDirectory(t *testing.T) {
	dir := HomeDirectory()
	if runtime.GOOS == "linux" {
		if !strings.Contains(dir, "/home/") {
			t.Errorf("Error Home directory.\nExpected: /home/ \nGot: %s", dir)
		}
	}
}

func TestBlocksDirectory(t *testing.T) {
	dir := BlocksDirectory()
	if runtime.GOOS == "linux" {
		if dir != HomeDirectory()+"/.viacoin/blocks" {
			t.Errorf("Error Blocks directory.\nExpected: %s\nGot: %s", HomeDirectory()+".viacoin/blocks", dir)
		}
	}
}

func TestBlockFiles(t *testing.T) {
	fmt.Println(BlockFiles()[0])
	file := BlockFiles()[0]
	if file != BlocksDirectory()+"/blk00000.dat" {
		t.Errorf("Error Blockfile.\nExpected: %s\nGot: %s", file, BlocksDirectory()+"/blk00000.dat")
	}
}

func TestBlocksReader(t *testing.T) {
	BlocksReader()
}