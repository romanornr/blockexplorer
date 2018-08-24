package Reorg

import (
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/romanornr/cyberchain/blockdata"
	"github.com/romanornr/cyberchain/database"
	"testing"
	"github.com/btcsuite/btcd/btcjson"
)

func BuildMockDatabase() {
	for i := int64(1); i < 6; i++ {
		blockhash := blockdata.GetBlockHash(i)
		block, _ := blockdata.GetBlock(blockhash)
		database.AddBlock(db, blockhash.String(), block)
		database.AddIndexBlockHeightWithBlockHash(db, blockhash.String(), block.Height)
	}
}

var hashes = [7]string{"ded7508b6b6452bfc99961366e3206a7a258cf897d3148b46e590bbf6f23f3d9",
	"45c2eb3f3ca602e36b9fac0c540cf2756f1d41719b4be25adb013f87bafee7bc",
	"7539b2ae01fd492adcc16c2dd8747c1562a702f9057560fee9ca647b67b729e2",
	"a35d1bdbd41ea6c290d9a151bdafd39b76eda3c9c9d44e02d0209dd77f5aec1f",
	"d8c9053f3c807b1465bd0a8bc99421e294066dd59e98cf14bb49d990ea88aff6", //5
	"e8957dac3477849c431dce6929e45ca829598bf45f05f776742f04f06c246ae7",
	"5ca78b039ccfec56373a4392c043bb9a6c77f8c2934af96b036c00dd2e4a0cfa",
}

// this block has a blockheight:2 but this blockheight is already in the database
// It has a different hash compared to the block in the database with blockheight:2
// This would be a potential chain reorg
var reorgBlock = &btcjson.GetBlockVerboseResult{
	Hash: "5ca78b039ccfec56373a4392c043bb9a6c77f8c2934af96b036c00dd2e4a0cfa",
	Height: 2,
}

func TestComparePreviousHash(t *testing.T) {

	BuildMockDatabase()

	//check if no chain reorg gets detected
	blockhash, err := chainhash.NewHashFromStr(hashes[5]) // check if blockhash is valid
	if err != nil {
		t.Errorf("The hash %s is not valid\n", blockhash.String())
	}

	block, err := blockdata.GetBlock(blockhash)
	if err != nil {
		t.Errorf("Could not get block: %s via RPC", hashes[5])
	}

	err = Check(block)

	if err != nil {
		t.Errorf("Did not expect reorg: %s", err)
	}

	//check if it detects a chain reorg
	//blockhash, _ = chainhash.NewHashFromStr(hashes[2]) // check if blockhash is valid
	//block, _ = blockdata.GetBlock(blockhash)
	err = Check(reorgBlock)

	if err == nil {
		t.Errorf("No chain reorg detected, however it was exected %s", err)
	}
}
