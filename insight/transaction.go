package insight

import "github.com/btcsuite/btcd/btcjson"

// insight api compatible address index struct
type AddrIndex struct {
	AddrStr                 string   `json:"addrStr"`
	Balance                 float64  `json:"balance"`
	BalanceSat              float64  `json:"balanceSat"`
	TotalReceived           float64  `json:"totalReceived"`
	TotalReceivedSat        float64  `json:"totalReceivedSat"`
	TotalSent               float64  `json:"totalSent"`
	TotalSentSat            float64  `json:"totalSentSat"`
	UnconfirmedBalance      float64  `json:"unconfirmedBalance"`
	UnconfirmedBalanceSat   float64  `json:"unconfirmedBalanceSat"`
	UnconfirmedTxApperances float64  `json:"unconfirmedTxApperances"`
	TxApperances            uint64   `json:"txApperances"`
	Transactions            []string `json:"transactions"`
}

// compatible insight tx struct
// MISSING: inside tge btcjson.TxRawResult there's []vin and []vout
// which also need exta fields
type TxRawResult struct {
	*btcjson.TxRawResult
	ValueOut float64 `json:"valueOut"`
	ValueIn  float64 `json:"valueIn"`
	Fees     float64 `json:"fees"`
}

// How []vin should look like. Address, ValueSatoshi and DoubleSpentTxID need
// to be added to the vin struct
//type Vin struct {
//	Coinbase        string      `json:"coinbase"`
//	Txid            string      `json:"txid"`
//	Vout            uint32      `json:"vout"`
//	ScriptSig       *ScriptSig  `json:"scriptSig"`
//	Sequence        uint32      `json:"sequence"`
//	Witness         []string    `json:"txinwitness"`
//	Address         string      `json:"addr"`
//	ValueSatoshi    int64       `json:"valueSat"`
//	DoubleSpentTxID interface{} `json:"doubleSpentTxID"`
//}

