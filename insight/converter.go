package insight

import (
	"github.com/btcsuite/btcd/btcjson"
	"github.com/romanornr/cyberchain/insightjson"
	"github.com/btcsuite/btcutil"
)

func ConvertToInsightBlock(block *btcjson.GetBlockVerboseResult) (*insightjson.BlockResult, error) {

	insightBlock := insightjson.BlockResult{
		Hash:              block.Hash,
		Size:              block.Size,
		Height:            block.Height,
		Version:           block.Version,
		MerkleRoot:        block.MerkleRoot,
		Tx:                block.Tx,
		Time:              block.Time,
		Nonce:             block.Nonce,
		Bits:              block.Bits,
		Difficulty:        block.Difficulty,
		Confirmations:     block.Confirmations,
		PreviousBlockHash: block.PreviousHash,
		NextBlockHash:     block.NextHash,
		IsMainChain:       true,
	}

	return &insightBlock, nil

}

func TxConverter(tx *btcjson.TxRawResult) ([]insightjson.Tx) {
	return ConvertToInsightTransaction(tx, false, false, false)
}

//NOTE: Address retrieval and vin needs to be fixed
func ConvertToInsightTransaction(tx *btcjson.TxRawResult, noAsm, noScriptSig, noSpent bool) []insightjson.Tx {

	var newTransaction []insightjson.Tx

	// TODO Blockheight
	txNew := insightjson.Tx{
		Txid: tx.Txid,
		Version: tx.Version,
		Locktime: tx.LockTime,
		Blockhash: tx.BlockHash,
		Confirmations: tx.Confirmations,
		Time: tx.Time,
		Blocktime: tx.Blocktime,
		Size: uint32(len(tx.Hex) /2),
	}

	var vInSum, vOutSum float64

	for vinID, vin := range tx.Vin {

		insightVin := &insightjson.Vin{
			Txid:     vin.Txid,
			Vout:     vin.Vout,
			Sequence: vin.Sequence,
			N:        vinID,
			CoinBase: vin.Coinbase,
		}

		//scriptpubkey
		if !noScriptSig {
			insightVin.ScriptSig = new(insightjson.ScriptSig)
			if vin.ScriptSig != nil {
				if !noAsm {
					insightVin.ScriptSig.Asm = vin.ScriptSig.Asm
				}
				insightVin.ScriptSig.Hex = vin.ScriptSig.Hex
			}

		}

		// address retrieval

		amount, _ := btcutil.NewAmount(insightVin.Value)
		insightVin.ValueSat = int64(amount)

		vInSum += insightVin.Value
		txNew.Vins = append(txNew.Vins, insightVin)
	}


	for _, v := range tx.Vout {
		InsightVout := &insightjson.Vout{
			Value: v.Value,
			N:     v.N,
			ScriptPubKey: insightjson.ScriptPubKey{
				Addresses: v.ScriptPubKey.Addresses,
				Type:      v.ScriptPubKey.Type,
				Hex:       v.ScriptPubKey.Hex,
			},
		}

		if !noAsm {
			InsightVout.ScriptPubKey.Asm = v.ScriptPubKey.Asm
		}

		txNew.Vouts = append(txNew.Vouts, InsightVout)
		vOutSum += v.Value
	}

	amount, _ := btcutil.NewAmount(vOutSum)
	txNew.ValueOut = amount.ToBTC()

	amount, _ = btcutil.NewAmount(vInSum)
	txNew.ValueIn = amount.ToBTC()

	amount, _ = btcutil.NewAmount(txNew.ValueIn - txNew.ValueOut)
	txNew.Fees = amount.ToBTC()

	if txNew.Vins != nil && txNew.Vins[0].CoinBase != "" {
		txNew.IsCoinBase = true
		txNew.ValueIn = 0
		txNew.Fees = 0
		for _, v := range txNew.Vins {
			v.Value = 0
			v.ValueSat = 0
		}
	}

	if !noSpent {
		//todo
	}

	newTransaction = append(newTransaction, txNew)
	return newTransaction
}
