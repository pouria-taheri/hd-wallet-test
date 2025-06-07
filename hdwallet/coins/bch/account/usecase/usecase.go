package usecase

import (
	"fmt"
	ad "git.mazdax.tech/blockchain/hdwallet/account/domain"
	"git.mazdax.tech/blockchain/hdwallet/coins/bch/account/domain"
	bd "git.mazdax.tech/blockchain/hdwallet/coins/bch/domain"
	"git.mazdax.tech/blockchain/hdwallet/config"
	"git.mazdax.tech/data-layer/configcore"
	"git.mazdax.tech/data-layer/loggercore/logger"
	"github.com/bchsuite/bchd/bchec"
	"github.com/bchsuite/bchd/chaincfg"
	"github.com/bchsuite/bchutil"
	"github.com/bchsuite/bchutil/hdkeychain"
	"strings"
)

type useCase struct {
	Config      domain.Config
	ChainParams *chaincfg.Params

	AddrMgr *Manager

	RootKey    *hdkeychain.ExtendedKey
	RootPubKey *hdkeychain.ExtendedKey

	coin   string
	logger logger.Logger
}

func New(logger logger.Logger, registry configcore.Registry,
	secureConfig config.SecureConfig, coin string) domain.UseCase {
	uc := &useCase{
		logger: logger,
		coin:   coin,
	}
	if err := registry.Unmarshal(&uc.Config); err != nil {
		panic(err)
	}
	uc.Config.Initialize(secureConfig)

	switch strings.ToLower(uc.Config.ChainType) {
	case "mainnet":
		uc.ChainParams = &chaincfg.MainNetParams
	case "regtest", "testnet":
		uc.ChainParams = &chaincfg.RegressionNetParams
	case "testnet3":
		uc.ChainParams = &chaincfg.TestNet3Params
	case "simnet":
		uc.ChainParams = &chaincfg.SimNetParams
	}
	return uc
}

func (uc *useCase) Initialize(scopes []ad.KeyScope, scopeAddrMap map[ad.KeyScope]ad.ScopeAddrSchema) {
	seed, err := uc.Config.GetSeedWithErrorChecking()
	if err != nil {
		panic(err)
	}

	// Generate the BIP0044 HD key structure to ensure the
	// provided seed can generate the required structure with no
	// issues.

	// Derive the master extended key from the seed.
	rootKey, err := hdkeychain.NewMaster(seed, uc.ChainParams)
	if err != nil {
		uc.logger.With(logger.Field{
			"submodule": "accountInitialize",
			"error":     err,
		}).ErrorF("cannot create masterKey from seed")
		panic("cannot create masterKey from seed")
	}
	rootPubKey, err := rootKey.Neuter()
	if err != nil {
		uc.logger.With(logger.Field{
			"submodule": "accountInitialize",
			"error":     err,
		}).ErrorF("failed to neuter master extended key")
		panic("failed to neuter master extended key")
	}
	uc.RootKey = rootKey
	uc.RootPubKey = rootPubKey

	uc.AddrMgr, err = LoadManager(uc.RootKey, uc.ChainParams, scopes, scopeAddrMap, uc.logger)

	if err != nil {
		uc.logger.With(logger.Field{
			"submodule": "accountInitialize",
			"error":     err,
		}).ErrorF("cannot load scope manager")
		panic("cannot load scope manager from key and scopes")
	}
}

func (uc *useCase) Coin() string {
	return uc.coin
}

func (uc *useCase) GetAccount(request ad.Request) (*ad.Account, error) {
	acc, err := uc.AddrMgr.GetAccount(request.KeyScope, request)
	return acc, err
}

func (uc *useCase) GetManagedAddress(request ad.Request) (ad.ManagedAddress, error) {
	managerAddr, err := uc.AddrMgr.GetManagedAddress(request.KeyScope, request)
	return managerAddr, err
}

func (uc *useCase) GetChainType() string {
	return uc.Config.ChainType
}

func (uc *useCase) SetChainParams(params *chaincfg.Params) {
	uc.ChainParams = params
}

func (uc *useCase) GetChainParams() *chaincfg.Params {
	return uc.ChainParams
}

func (uc *useCase) GetKey(address bchutil.Address) (*bchec.PrivateKey, bool, error) {
	addr := address.(bd.AddressModel)
	dp := addr.GetDerivationPath()
	ks := addr.GetKeyScope()
	ma, err := uc.GetManagedAddress(ad.Request{
		DerivationPath: ad.DerivationPath{
			Account: dp.Account,
			Branch:  dp.Branch,
			Index:   dp.Index,
		},
		KeyScope: ad.KeyScope{
			Purpose: ks.Purpose,
			Coin:    ks.Coin,
		},
		Private: addr.IsPrivate(),
	})
	if err != nil {
		return nil, false, err
	}
	mpka, ok := ma.(domain.ManagedPubKeyAddress)
	if !ok {
		uc.logger.With(logger.Field{
			"submodule": "account usecase",
			"section":   "get key",
			"address":   addr,
			"type":      ma,
		}).ErrorF("managed address type waddrmgr.ManagedPubKeyAddress expected but got...")
		e := fmt.Errorf("managed address type is not valid")
		return nil, false, e
	}
	privKey, err := mpka.PrivKey()
	if err != nil {
		uc.logger.With(logger.Field{
			"submodule": "account usecase",
			"section":   "get key",
			"type":      ma,
		}).ErrorF("error in get private key from address")
		return nil, false, err
	}
	return privKey, ma.Compressed(), nil
}

func (uc *useCase) GetScript(address bchutil.Address) ([]byte, error) {
	addr := address.(bd.AddressModel)
	dp := addr.GetDerivationPath()
	ks := addr.GetKeyScope()
	ma, err := uc.GetManagedAddress(ad.Request{
		DerivationPath: ad.DerivationPath{
			Account: dp.Account,
			Branch:  dp.Branch,
			Index:   dp.Index,
		},
		KeyScope: ad.KeyScope{
			Purpose: ks.Purpose,
			Coin:    ks.Coin,
		},
		Private: addr.IsPrivate(),
	})
	if err != nil {
		uc.logger.With(logger.Field{
			"submodule": "account usecase",
			"section":   "get script",
			"address":   addr,
			"error":     err,
		}).ErrorF("error in get managed address")
		return nil, err
	}

	msa, ok := ma.(domain.ManagedScriptAddress)
	if !ok {
		uc.logger.With(logger.Field{
			"submodule": "account usecase",
			"section":   "get script",
			"address":   addr,
			"type":      ma,
		}).ErrorF("managed address type waddrmgr.ManagedScriptAddress expected but got...")
		e := fmt.Errorf("managed address type is not valid")
		return nil, e
	}
	return msa.Script()
}
