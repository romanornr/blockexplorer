package testdata

import (
	"github.com/romanornr/cyberchain/insightjson"
)

// Test transactions
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

// Test transactions
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
