package subsidy

import (
	"github.com/viacoin/viad/chaincfg"
	"github.com/btcsuite/btcutil"
)

const baseSubsidy = 50 * btcutil.SatoshiPerBitcoin

func CalcViacoinBlockSubsidy(height int32, chainParams *chaincfg.Params) int64{
	if chainParams.GenerateSupported { //regtest: use bitcoin schedule
		if chainParams.SubsidyReductionInterval == 0 {
			return baseSubsidy
		}
		// Equivalent to: baseSubsidy / 2^(height/subsidyHalvingInterval)
		return baseSubsidy >> uint(height/chainParams.SubsidyReductionInterval)
	}

	// Viacoin schedule
	var zeroRewardHeight int32
	if chainParams.ReduceMinDifficulty {
		zeroRewardHeight = 2001
	} else {
		zeroRewardHeight = 10001
	}
	rampHeight := zeroRewardHeight + 43200 //4 periods of 10800

	subsidy := int64(0)
	if height == 0 {
		subsidy = 0
	} else if height == 1 {
		subsidy = 10000000 * btcutil.SatoshiPerBitcoin
	} else if height <= zeroRewardHeight + 10800 {
		// first 10800 block after zero reward period is 10 coins per block
		subsidy = 10 * btcutil.SatoshiPerBitcoin
	} else if height <= rampHeight {
		// every 10800 blocks reduce nSubsidy from 8 to 6
		subsidy = (8 - int64((height - zeroRewardHeight - 1) / 10800)) * btcutil.SatoshiPerBitcoin
	} else if height <= 1971000 {
		subsidy = 5 * btcutil.SatoshiPerBitcoin
	} else {
		halvings := uint32(height / chainParams.SubsidyReductionInterval)
		if halvings <= 64 {
			subsidy = 20 * btcutil.SatoshiPerBitcoin
			subsidy >>= halvings
		}
	}
	return subsidy
}
