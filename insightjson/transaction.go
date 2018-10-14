package insightjson

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

type Vin struct {
	Txid      string     `json:"txid,omitempty"`
	Vout      uint32     `json:"vout,omitempty"`
	Sequence  uint32     `json:"sequence,omitempty"`
	N         int        `json:"n"`
	ScriptSig *ScriptSig `json:"scriptSig,omitempty"`
	Addr      string     `json:"addr,omitempty"`
	ValueSat  int64      `json:"valueSat"`
	Value     float64    `json:"value,omitempty"`
	Coinbase  string     `json:"coinbase,omitempty"`
}

type Vout struct {
	Value        float64             `json:"value"`
	N            uint32              `json:"n"`
	ScriptPubKey ScriptPubKey `json:"scriptPubKey,omitempty"`
	SpentTxID    interface{}         `json:"spentTxId"`
	SpentIndex   interface{}         `json:"spentIndex"`
	SpentHeight  interface{}         `json:"spentHeight"`
}

type ScriptPubKey struct {
	Hex       string   `json:"hex,omitempty"`
	Asm       string   `json:"asm,omitempty"`
	Addresses []string `json:"addresses,omitempty"`
	Type      string   `json:"type,omitempty"`
}

type ScriptSig struct {
	Hex string `json:"hex,omitempty"`
	Asm string `json:"asm,omitempty"`
}