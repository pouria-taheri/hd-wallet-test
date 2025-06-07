package domain

import (
	ad "git.mazdax.tech/blockchain/hdwallet/account/domain"
)

type CardanoWalletModel interface {
	Coin() string
	GetAccount(request ad.Request) (*ad.Account, error)
	GetAddresses() ([]byte, error)
	GetWallet() ([]byte, error)
	NewAddress() ([]byte, error)
}
