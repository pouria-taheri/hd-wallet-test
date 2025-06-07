package wallet

import (
	"git.mazdax.tech/blockchain/hdwallet/coins/bsc/domain"
	"git.mazdax.tech/data-layer/loggercore/logger"
	"github.com/btcsuite/btcd/chaincfg"
	"strings"
)

func NewBSCHdWallet(config *domain.BscConfig, logger logger.Logger) domain.BSCWallet {
	hd := new(bscWallet)
	hd.MasterKey = createMasterKey(config.SecureConfig, getChainParam(config.ChainType), logger)
	hd.WalletMainAcc = createWalletMainAccount(hd.MasterKey.MasterKey, logger)
	return *hd
}

func getChainParam(chainType string) (ChainParams *chaincfg.Params) {
	switch strings.ToLower(chainType) {
	case "mainnet":
		ChainParams = &chaincfg.MainNetParams
	case "regression":
		ChainParams = &chaincfg.RegressionNetParams
	case "testnet3", "testnet":
		ChainParams = &chaincfg.TestNet3Params
	case "simnet":
		ChainParams = &chaincfg.SimNetParams
	}
	return ChainParams
}
