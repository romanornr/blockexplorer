package insight

import (
	"fmt"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/romanornr/blockexplorer/blockdata"
	"testing"
)

func TestConvertToInsightTransaction(t *testing.T) {
	x := blockdata.GetBlockHash(2)
	xx, _ := blockdata.GetBlock(x)

	hash, _ := chainhash.NewHashFromStr(xx.Tx[0])

	tx := blockdata.GetRawTransactionVerbose(hash)
	result := ConvertToInsightTransaction(tx, false, false, false)
	fmt.Println(result)
}

var rawTx = btcjson.TxRawResult{
	Hex:      "",
	Txid:     "583910b7bf90ab802e22e5c25a89b59862b20c8c1aeb24dfb94e7a508a70f121",
	Hash:     "583910b7bf90ab802e22e5c25a89b59862b20c8c1aeb24dfb94e7a508a70f121",
	Size:     225,
	Vsize:    225,
	Version:  1,
	LockTime: 0,
}
