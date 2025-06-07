package wallet

import (
	"fmt"
	"git.mazdax.tech/blockchain/hdwallet/coins/xlm/domain"
	"git.mazdax.tech/data-layer/loggercore/logger"
	"github.com/aliworkshop/stellar-go/exp/crypto/derivation"
	"github.com/aliworkshop/stellar-go/keypair"
)

func createWalletMainAccount(key domain.MasterKey, log logger.Logger) domain.WalletMainAcc {
	wma := domain.WalletMainAcc{}

	pair, err := keypair.FromRawSeed(key.MasterKey.RawSeed())
	if err != nil {
		log.With(logger.Field{
			"submodule": "walletMainAccount",
			"error":     err,
		}).FatalF("error in get main account from masterKey")
	}

	wma.PublicKey, _ = key.MasterKey.PublicKey()

	wma.Account = pair
	log.With(logger.Field{
		"submodule": "walletMainAccount",
		"address":   wma.Account.Address(),
	}).InfoF("Main account")

	return wma
}

func pathFromUserID(uid uint32) string {
	return fmt.Sprintf(derivation.StellarAccountPathFormat, uid)
}
