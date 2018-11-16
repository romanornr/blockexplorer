package testdata

import (
	"github.com/romanornr/cyberchain/insightjson"
)

// Test transactions.
var (
	TxBlock5421176Insight = &insightjson.Tx{
		Txid:     "31c0cbc8411de76eac6018183e96d1cc2c904a9b50096758041eec92d9c9b9f9",
		Version:  2,
		Locktime: 5421173,
		Vins: []*insightjson.Vin{
			{
				Txid:     "d144a6928043424ca6cde94491f1b642bc976471d1d1b103b592c49903d1544b",
				Sequence: 4294967294,
				N:        0,
				ScriptSig: &insightjson.ScriptSig{
					Hex: "473044022002911635530b9f5a38af1d3b30021b2ba4c764e0bfe5eb51ac64d2226f0fc9e602200f65b56ca631abf8e21d29024fd06f2d65292c9819b8674c402f0ff5a4827aed0121037e008989be991f383b39316deada461abaaf748110e5bcabe5e81e84509d1c8c",
					Asm: "3044022002911635530b9f5a38af1d3b30021b2ba4c764e0bfe5eb51ac64d2226f0fc9e602200f65b56ca631abf8e21d29024fd06f2d65292c9819b8674c402f0ff5a4827aed[ALL] 037e008989be991f383b39316deada461abaaf748110e5bcabe5e81e84509d1c8c",
				},
				ValueSat: 297954800,
				Value:    2.979548,
				Addr:     "Vxs7pdwf4XWt1MMKedMhewH9dCbwonFLav",
				Vout:     0,
				// DoubleSpentTxID: nil,
			},
		},
		Vouts: []*insightjson.Vout{
			{
				Value: 2,
				N:     0,
				ScriptPubKey: insightjson.ScriptPubKey{
					Hex:       "76a91456c7359ed52d61c1ca371d7dc136632148169c5e88ac",
					Asm:       "OP_DUP OP_HASH160 56c7359ed52d61c1ca371d7dc136632148169c5e OP_EQUALVERIFY OP_CHECKSIG",
					Addresses: []string{"VhuffXKNA3j9hgp2JYGrj6uHQ6KUU6zNbS"},
					Type:      "pubkeyhash",
				},
				SpentTxID:   "d78999b2ad131bd393c06738bd34996da80a556d6b1e9486447a023b91ef6ea3",
				SpentIndex:  0,
				SpentHeight: 5422075,
			},
			{
				Value: 0.979096,
				N:     1,
				ScriptPubKey: insightjson.ScriptPubKey{
					Hex:       "76a9147fbf8dfb4c104984c1824dc1c129a1f2bd6ea91b88ac",
					Asm:       "OP_DUP OP_HASH160 7fbf8dfb4c104984c1824dc1c129a1f2bd6ea91b OP_EQUALVERIFY OP_CHECKSIG",
					Addresses: []string{"VmeJCXAxkR5LxEwezdsGKtNoxet8A63VVX"},
					Type:      "pubkeyhash",
				},
				SpentTxID:   "34e336269c45be83d6892379258844e5508380d87cee4533e4404471c106c783",
				SpentIndex:  3,
				SpentHeight: 5421176,
			},
		},
		Blockhash:     "0d37d5dedab84e4c70a35113acbbf2c3514a46e66e6ff1aaae9b2ece846a3e63",
		Blockheight:   5421176,
		Confirmations: 207783, // should we even have confirmations? We can always count that using blocks.
		Time:          1536613229,
		Blocktime:     1536613229,
		ValueOut:      2.979096,
		Size:          225,
		ValueIn:       2.979548,
		Fees:          0.000452,
	}
)

// Test addresses.
var (
	AddrStdVlav = &insightjson.Address{
		Address: "Vxs7pdwf4XWt1MMKedMhewH9dCbwonFLav",
	}

	AddrStdVLavInfo = &insightjson.AddressInfo{
		Address:               "Vxs7pdwf4XWt1MMKedMhewH9dCbwonFLav",
		Balance:               0,
		BalanceSat:            0,
		TotalReceived:         2.979548,
		TotalReceivedSat:      297954800,
		TotalSent:             2.979548,
		TotalSentSat:          297954800,
		UnconfirmedBalance:    0,
		UnconfirmedBalanceSat: 0,
		TxAppearances:         2,
		TransactionsID: []string{
			"31c0cbc8411de76eac6018183e96d1cc2c904a9b50096758041eec92d9c9b9f9",
			"d144a6928043424ca6cde94491f1b642bc976471d1d1b103b592c49903d1544b",
		},
	}
)

// Test block.
var (
	Block2 = &insightjson.BlockResult{
		Hash:              "45c2eb3f3ca602e36b9fac0c540cf2756f1d41719b4be25adb013f87bafee7bc",
		Size:              202,
		Height:            2,
		Version:           2,
		MerkleRoot:        "bf5721dfb2a84b8f39ec28dd5a7d4e8b219ca3a361dd44db2d67470818a625ea",
		Tx:                []string{"bf5721dfb2a84b8f39ec28dd5a7d4e8b219ca3a361dd44db2d67470818a625ea"},
		Time:              1405608158,
		Nonce:             554156041,
		Bits:              "1e01ffff",
		Difficulty:        0.001953110098725118,
		Confirmations:     257,
		PreviousBlockHash: "5ca83af67146e286610e118cc8f8e6a183c319fbb4a8fdb9e99daa2b8a29b3e3",
		NextBlockHash:     "7539b2ae01fd492adcc16c2dd8747c1562a702f9057560fee9ca647b67b729e2",
		IsMainChain:       true,
	}

	Block3 = &insightjson.BlockResult{
		Hash:              "7539b2ae01fd492adcc16c2dd8747c1562a702f9057560fee9ca647b67b729e2",
		Size:              202,
		Height:            3,
		Version:           2,
		MerkleRoot:        "cef916ad6fc1c3ca4ea50360f68ff0a43b5b4ffc217a51c1128106a61ced9900",
		Tx:                []string{"cef916ad6fc1c3ca4ea50360f68ff0a43b5b4ffc217a51c1128106a61ced9900"},
		Time:              1405608158,
		Nonce:             191037457,
		Bits:              "1e01ffff",
		Difficulty:        1953.11,
		Confirmations:     5629692,
		PreviousBlockHash: "45c2eb3f3ca602e36b9fac0c540cf2756f1d41719b4be25adb013f87bafee7bc",
		NextBlockHash:     "a35d1bdbd41ea6c290d9a151bdafd39b76eda3c9c9d44e02d0209dd77f5aec1f",
		IsMainChain:       true,
	}
)