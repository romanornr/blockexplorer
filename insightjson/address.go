// Copyright (c) 2018-2019, The Decred developers
// Copyright (c) 2017, Jonathan Chappelow
// Copyright (c) 2019, Romano, Viacoin developer
// See LICENSE for details.

package insightjson

type Address struct {
	Address      string `json:"address, omitempty"`
	From         int    `json:"from"`
	To           int    `json:"to"`
	Transactions []Tx   `json:"transactions, omitempty"`
}

type AddressInfo struct {
	Address                  string   `bson:"addrStr,omitempty" json:"addrStr,omitempty"`
	Balance                  float64  `bson:"balance" json:"balance"`
	BalanceSat               int64    `bson:"balanceSat" json:"balanceSat"`
	TotalReceived            float64  `bson:"totalReceived" json:"totalReceived"`
	TotalReceivedSat         int64    `bson:"totalReceivedSat" json:"totalReceivedSat"`
	TotalSent                float64  `bson:"totalSent" json:"totalSent"`
	TotalSentSat             int64    `bson:"totalSentSat" json:"totalSentSat"`
	UnconfirmedBalance       float64  `bson:"unconfirmedBalance" json:"unconfirmedBalance"`
	UnconfirmedBalanceSat    int64    `bson:"unconfirmedBalanceSat" json:"unconfirmedBalanceSat"`
	UnconfirmedTxAppearances int64    `bson:"unconfirmedTxAppearances" json:"unconfirmedTxAppearances"`
	TxAppearances            int64    `bson:"txAppearances" json:"txAppearances"`
	TransactionsID           []string `bson:"transactions,omitempty" json:"transactions,omitempty"`
}

// address tx output
type AddressTxnOutput struct {
	Address       string  `json:"address"`
	TxnID         string  `json:"txid"`
	Vout          uint32  `json:"vout"`
	Blocktime     int64   `json:"ts,omitempty"`
	ScriptPubKey  string  `json:"scriptPubKey"`
	Height        int64   `json:"height,omitempty"`
	BlockHash     string  `json:"block_hash,omitempty"`
	Amount        float64 `json:"amount,omitempty"`
	Atoms         int64   `json:"atoms,omitempty"`
	Satoshis      int64   `json:"satoshis,omitempty"`
	Confirmations int64   `json:"confirmations"`
}

// return from GetSpendDetailsByFundingHash
type SendByFundingHash struct {
	FundingTxVoutIndex uint32
	SpendingTxVinIndex interface{}
	SpendingTxHash     interface{}
	BlockHeight        interface{}
}
