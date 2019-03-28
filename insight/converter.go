// Copyright (c) 2018, The Decred developers
// Copyright (c) 2017, The dcrdata developers
// Copyright (c) 2019, Romano, Viacoin developer
// See LICENSE for details.

package insight

import (
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcutil"
	"github.com/romanornr/cyberchain/insightjson"
	"github.com/romanornr/cyberchain/mongodb"
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

func TxConverter(tx *btcjson.TxRawResult, blockheight int64) []insightjson.Tx {
	return ConvertToInsightTransaction(tx, blockheight, false, false, false)
}

// NOTE: Address retrieval and vin needs to be fixed
func ConvertToInsightTransaction(tx *btcjson.TxRawResult, blockheight int64, noAsm, noScriptSig, noSpent bool) []insightjson.Tx {

	var newTransaction []insightjson.Tx

	// TODO Blockheight
	txNew := insightjson.Tx{
		Txid:          tx.Txid,
		Version:       tx.Version,
		Locktime:      tx.LockTime,
		Blockhash:     tx.BlockHash,
		Blockheight:   blockheight,
		Confirmations: tx.Confirmations,
		Time:          tx.Time,
		Blocktime:     tx.Blocktime,
		Size:          uint32(len(tx.Hex) / 2),
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

		// scriptpubkey
		if !noScriptSig {
			insightVin.ScriptSig = new(insightjson.ScriptSig)
			if vin.ScriptSig != nil {
				if !noAsm {
					insightVin.ScriptSig.Asm = vin.ScriptSig.Asm
				}
				insightVin.ScriptSig.Hex = vin.ScriptSig.Hex
			}
		}

		// retrieval for vin[] to get addr and value // TODO What if there are multiple vins?
		vinHash, _ := chainhash.NewHashFromStr(vin.Txid)
		vinDbTx, err := mongodb.GetTransaction(*vinHash)
		if err == nil {
			if tx.Confirmations != 0 {
				i := insightVin.Vout
				insightVin.Value = vinDbTx.Vouts[i].Value
				insightVin.Addr = vinDbTx.Vouts[i].ScriptPubKey.Addresses[0]
			}
		}

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

	//nospent
	if !noSpent {
		// check https://explorer.viacoin.org/api/tx/99b5fa6e08c2d319dd68a16f781e72903a8e686f1c64dc1d11040490fbe81320 (call this txA)
		// Check the Vout and you see a spentTxID 7ec15e6386f019f488cc8ea418be32ae968335b27e44fa2ba37c20cfebc56ab2 (call this txB)
		// This is a transaction that happened after the current transaction that got processed.
		// current tx (txA) from explorer.viacoin has blockheight 3673
		// and txB happened at block 3674.

		//So when we process txB we need to "update" txA
		//txB->vin[0] has txA->txid as value
		//txA->vout[0].SpentTxId has txB->txid
		//txA->Spentheight has txB->blockheight
		for _, vin := range txNew.Vins {
			if len(vin.Txid) < 1 {
				continue
			}

			txHash, _ := chainhash.NewHashFromStr(vin.Txid)
			tx, err := mongodb.GetTransaction(*txHash)
			if err == nil {
				i := vin.Vout
				tx.Vouts[i].SpentTxID = txNew.Txid

				for idx, x := range txNew.Vins {
					if tx.Txid == x.Txid {
						tx.Vouts[i].SpentIndex = idx
					}
				}

				tx.Vouts[i].SpentHeight = txNew.Blockheight
				go mongodb.UpdateTransaction(&tx)
			}
		}
	}
	// If you are reading this all this, you have no fucking life.

	newTransaction = append(newTransaction, txNew)
	return newTransaction
}
