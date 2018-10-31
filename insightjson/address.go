package insightjson

type Address struct {
	Address      string `json:"address, omitempty"`
	From         int    `json:"from"`
	To           int    `json:"to"`
	Transactions []Tx   `json:"transactions, omitempty"`
}

type AddressInfo struct {
	Address                  string   `json:"addrStr,omitempty"`
	Balance                  float64  `json:"balance"`
	BalanceSat               int64    `json:"balanceSat"`
	TotalReceived            float64  `json:"totalReceived"`
	TotalReceivedSat         int64    `json:"totalReceivedSat"`
	TotalSent                float64  `json:"totalSent"`
	TotalSentSat             float64  `json:"totalSentSat"`
	UnconfirmedBalance       float64  `json:"unconfirmedBalance"`
	UnconfirmedBalanceSat    float64  `json:"unconfirmedBalanceSat"`
	UnconfirmedTxAppearances int64    `json:"unconfirmedTxApperances"`
	TxAppearances            int64    `json:"txApperances"`
	TransactionsID           []string `json:"transactions,omitempty"`
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
	SpendingTxHash interface{}
	BlockHeight interface{}
}