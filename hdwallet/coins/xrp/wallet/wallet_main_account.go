package wallet

import (
	"fmt"
	"git.mazdax.tech/blockchain/hdwallet/coins/xrp/crypto"
	"git.mazdax.tech/blockchain/hdwallet/coins/xrp/domain"
	"git.mazdax.tech/data-layer/loggercore/logger"
)

func createWalletMainAccount(key domain.MasterKey, log logger.Logger) domain.WalletMainAcc {
	wma := domain.WalletMainAcc{}

	public, errPub := key.MasterKey.PublicKey()
	if errPub != nil {
		log.With(logger.Field{
			"submodule": "masterKey",
			"error":     errPub,
		}).FatalF("cannot create public key from private key")
	}

	xrpKey, errPrv := crypto.NewEd25519KeyFromPrivate(append(key.MasterKey.Key, public...))
	if errPrv != nil {
		log.With(logger.Field{
			"submodule": "masterKey",
			"error":     errPrv,
		}).FatalF("cannot create ed25519 key from private key")
	}
	sourceAddress, err := crypto.NewAccountId(xrpKey.Id(nil))
	if err != nil {
		log.With(logger.Field{
			"submodule": "masterKey",
			"error":     err,
		}).FatalF("cannot create address from key")
	}

	wma.PublicKey = sourceAddress.String()

	log.With(logger.Field{
		"submodule": "walletMainAccount",
		"address":   wma.PublicKey,
	}).InfoF("Main account")

	return wma
}

func pathFromUserId(uid uint64, index uint32) string {
	return fmt.Sprintf("m/44'/144'/%d'/0'/%d'", uid, index)
}
