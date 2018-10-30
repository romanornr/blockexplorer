package insight

import (
	"testing"
	"github.com/romanornr/cyberchain/blockdata"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"fmt"
)

func TestConvertToInsightTransaction(t *testing.T) {
	x := blockdata.GetBlockHash(2)
	xx, _:= blockdata.GetBlock(x)

	hash,_ := chainhash.NewHashFromStr(xx.Tx[0])

	tx := blockdata.GetRawTransactionVerbose(hash)
	result := ConvertToInsightTransaction(tx)
	fmt.Println(result)
}
