package domain

import (
	"git.mazdax.tech/blockchain/hdwallet/config"
)

type Config struct {
	config.SecureConfig
	// represents coin type, e.g. bch, eth, etc.
	Coin string
	// represents type of chain
	ChainType string
}

func (c *Config) Initialize(secureConfig config.SecureConfig) {
	c.SecureConfig = secureConfig
}
