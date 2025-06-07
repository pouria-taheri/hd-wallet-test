package domain

import (
	"github.com/btcsuite/btcd/chaincfg"
	"strings"
)

type Config struct {
	// represents type of chain
	ChainType string

	ChainParams *chaincfg.Params
}

func (c *Config) Initialize() {
	switch strings.ToLower(c.ChainType) {
	case "mainnet":
		c.ChainParams = &chaincfg.MainNetParams
	case "regression":
		c.ChainParams = &chaincfg.RegressionNetParams
	case "testnet3":
		c.ChainParams = &chaincfg.TestNet3Params
	case "simnet":
		c.ChainParams = &chaincfg.SimNetParams
	}
}
