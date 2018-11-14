package subsidy

import (
	"github.com/viacoin/viad/chaincfg"
	"github.com/viacoin/viad/blockchain"
	"github.com/btcsuite/btcutil"
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