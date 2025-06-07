package domain

import (
	ad "git.mazdax.tech/blockchain/hdwallet/account/domain"
	bd "git.mazdax.tech/blockchain/hdwallet/coins/doge/domain"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
)

type UseCase interface {
	Coin() string
	Initialize(scopes []ad.KeyScope, scopeAddrMap map[ad.KeyScope]ad.ScopeAddrSchema)
	GetAccount(request ad.Request) (*ad.Account, error)
	GetManagedAddress(request ad.Request) (ad.ManagedAddress, error)

	GetChainType() string
	SetChainParams(params *chaincfg.Params)
	GetChainParams() *chaincfg.Params

	GetKey(addr *bd.Address) (*btcec.PrivateKey, bool, error)
	GetScript(addr *bd.Address) ([]byte, error)
}
