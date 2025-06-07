package domain

import (
	ad "git.mazdax.tech/blockchain/hdwallet/account/domain"
	"github.com/bchsuite/bchd/bchec"
	"github.com/bchsuite/bchd/chaincfg"
	"github.com/bchsuite/bchutil"
)

type UseCase interface {
	Coin() string
	Initialize(scopes []ad.KeyScope, scopeAddrMap map[ad.KeyScope]ad.ScopeAddrSchema)
	GetAccount(request ad.Request) (*ad.Account, error)
	GetManagedAddress(request ad.Request) (ad.ManagedAddress, error)

	GetChainType() string
	SetChainParams(params *chaincfg.Params)
	GetChainParams() *chaincfg.Params

	GetKey(addr bchutil.Address) (*bchec.PrivateKey, bool, error)
	GetScript(addr bchutil.Address) ([]byte, error)
}

