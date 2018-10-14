package insight

type Tx struct {
	Txid          string  `json:"txid,omitempty"`
	Version       int32   `json:"version,omitempty"`
	Locktime      uint32  `json:"locktime"`
	IsCoinBase    bool    `json:"isCoinBase,omitempty"`
	Vins          []*Vin  `json:"vin,omitempty"`
	Vouts         []*Vout `json:"vout,omitempty"`
	Blockhash     string  `json:"blockhash,omitempty"`
	Blockheight   int64   `json:"blockheight"`
	Confirmations int64   `json:"confirmations"`
	Time          int64   `json:"time,omitempty"`
	Blocktime     int64   `json:"blocktime,omitempty"`
	ValueOut      float64 `json:"valueOut,omitempty"`
	Size          uint32  `json:"size,omitempty"`
	ValueIn       float64 `json:"valueIn,omitempty"`
	Fees          float64 `json:"fees,omitempty"`
}
