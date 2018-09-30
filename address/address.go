package address

type Index struct {
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
