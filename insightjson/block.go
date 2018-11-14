package insightjson

type BlockResult struct {
	Hash              string   `json:"hash"`
	Size              int32    `json:"size"`
	Height            int64    `json:"height"`
	Version           int32    `json:"version"`
	MerkleRoot        string   `json:"merkleroot"`
	Tx                []string `json:"tx,omitempty"`
	Time              int64    `json:"time"`
	Nonce             uint32   `json:"nonce"`
	Bits              string   `json:"bits"`
	Difficulty        float64  `json:"difficulty"`
	Confirmations     int64    `json:"confirmations"`
	PreviousBlockHash string   `json:"previousblockhash"`
	NextBlockHash     string   `json:"nextblockhash,omitempty"`
	Reward            float64  `json:"reward"`
	IsMainChain       bool     `json:"isMainChain"`
	PoolInfo          *Pools   `json:"poolInfo"`
}

// TODO: implement block validation method
func (b *BlockResult) Validate() error {
	return nil
}

type Pools struct {
	PoolName      string   `json:"poolName"`
	URL           string   `json:"url,omitempty"`
	SearchStrings []string `json:"searchStrings,omitempty"`
}
