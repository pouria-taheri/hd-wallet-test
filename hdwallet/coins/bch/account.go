package bch

import (
	ad "git.mazdax.tech/blockchain/hdwallet/account/domain"
	"git.mazdax.tech/blockchain/hdwallet/coins/bch/account"
	"git.mazdax.tech/blockchain/hdwallet/config"
	"git.mazdax.tech/data-layer/configcore"
	"git.mazdax.tech/data-layer/loggercore/logger"
)

func NewAccountUseCase(registry configcore.Registry,
	secureConfig config.SecureConfig, logger logger.Logger) ad.UseCase {

	uc := account.NewUseCase(registry, secureConfig, "bch", logger)

	chainParams := uc.GetChainParams()
	var (
		// KeyScopeBIP0044 is the key scope for BIP0044 derivation. Legacy
		// wallets will only be able to use this key scope, and no keys beyond
		// it.
		KeyScopeBIP0044 = ad.KeyScope{
			Purpose: 44,
			Coin:    chainParams.HDCoinType,
		}

		// DefaultKeyScopes is the set of default key scopes that will be
		// created by the root manager upon initial creation.
		DefaultKeyScopes = []ad.KeyScope{
			KeyScopeBIP0044,
		}

		// ScopeAddrMap is a map from the default key scopes to the scope
		// address schema for each scope type. This will be consulted during
		// the initial creation of the root key manager.
		ScopeAddrMap = map[ad.KeyScope]ad.ScopeAddrSchema{
			KeyScopeBIP0044: {
				InternalAddrType: ad.PubKeyHash,
				ExternalAddrType: ad.PubKeyHash,
			},
		}
	)

	uc.Initialize(DefaultKeyScopes, ScopeAddrMap)
	return uc
}
