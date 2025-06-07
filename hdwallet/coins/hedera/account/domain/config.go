package domain

import (
	"git.mazdax.tech/blockchain/hdwallet/config"
)

type Config struct {
	config.SecureConfig
	Coin      string
	ChainType string
}

func (c *Config) Initialize(secureConfig config.SecureConfig) {
	c.SecureConfig = secureConfig
} 