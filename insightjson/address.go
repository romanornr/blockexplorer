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
