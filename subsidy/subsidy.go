package subsidy

import (
	"github.com/btcsuite/btcutil"
	"github.com/viacoin/viad/blockchain"
	"github.com/viacoin/viad/chaincfg"
)

// TODO Bitcoin, Litecoin & Other coins

// blocksubsidy for Viacoin
func CalcViacoinBlockSubsidy(height int32, isMainChain bool) float64 {
	var rewardInSatoshi int64

	if !isMainChain {
		rewardInSatoshi = blockchain.CalcBlockSubsidy(height, &chaincfg.TestNet3Params)
		return btcutil.Amount(rewardInSatoshi).ToBTC()
	}

	rewardInSatoshi = blockchain.CalcBlockSubsidy(height, &chaincfg.MainNetParams)
	return btcutil.Amount(rewardInSatoshi).ToBTC()
}
