package domain

import (
	bd "git.mazdax.tech/blockchain/hdwallet/coins/bch/domain"
)

type Config struct {
	bd.Config
}

func (c *Config) Initialize() {
	c.Config.Initialize()
}
