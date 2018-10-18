package insightjson

type BlockResult struct {
	Hash          string   `json:"hash"`
	Confirmations int64    `json:"confirmations"`
	Size          int32    `json:"size"`
	Height        int64    `json:"height"`
	Version       int32    `json:"version"`
	MerkleRoot    string   `json:"merkleroot"`
	Tx            []string `json:"tx,omitempty"`
	Time          int64    `json:"time"`
	Nonce         uint32   `json:"nonce"`
	Bits          string   `json:"bits"`
	Difficulty    float64  `json:"difficulty"`
	PreviousHash  string   `json:"previousblockhash"`
	NextHash      string   `json:"nextblockhash,omitempty"`
	Reward        float64  `json:"reward"`
	IsMainChain   bool     `json:"isMainChain"`
}
