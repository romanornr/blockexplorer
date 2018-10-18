package insight

import (
	"github.com/btcsuite/btcd/btcjson"
	"github.com/romanornr/cyberchain/insightjson"
)

func ConvertToInsightBlock(block *btcjson.GetBlockVerboseResult) (*insightjson.BlockResult, error) {

	insightBlock := insightjson.BlockResult{
		Hash: block.Hash,
		Size: block.Size,
		Height: block.Height,
		Version: block.Version,
		MerkleRoot: block.MerkleRoot,
		Tx: block.Tx,
		Time: block.Time,
		Nonce: block.Nonce,
		Bits: block.Bits,
		Difficulty: block.Difficulty,
		Confirmations: block.Confirmations,
		PreviousBlockHash: block.PreviousHash,
		NextBlockHash: block.NextHash,
		IsMainChain: true,
	}

	return &insightBlock, nil

}
