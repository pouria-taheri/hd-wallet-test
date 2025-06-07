package domain

import (
	bd "git.mazdax.tech/blockchain/hdwallet/coins/ltc/domain"
)

type Config struct {
	bd.Config
}

func (c *Config) Initialize() {
	c.Config.Initialize()
}
