package domain

import (
	bd "git.mazdax.tech/blockchain/hdwallet/coins/doge/domain"
)

type Config struct {
	bd.Config
}

func (c *Config) Initialize() {
	c.Config.Initialize()
}
