package wallet

import (
	"fmt"
	"git.mazdax.tech/blockchain/hdwallet/coins/eos/domain"
	"git.mazdax.tech/data-layer/loggercore/logger"
	"github.com/eoscanada/eos-go/btcsuite/btcutil/base58"
	"github.com/eoscanada/eos-go/ecc"
)

func createWalletMainAccount(key domain.MasterKey, log logger.Logger) domain.WalletMainAcc {
	wma := domain.WalletMainAcc{}

	wif := base58.CheckEncode(key.MasterKey.Key, '\x80')
	privateKey, errPrv := ecc.NewPrivateKey(wif)
	if errPrv != nil {
		log.With(logger.Field{
			"submodule": "masterKey",
			"error":     errPrv,
		}).FatalF("cannot create private key from wif")
	}

	wma.PrivateKey = privateKey
	wma.PublicKey = wma.PrivateKey.PublicKey()

	log.With(logger.Field{
		"submodule": "walletMainAccount",
		"address":   wma.PublicKey.String(),
	}).InfoF("Main account")

	return wma
}

func pathFromUserId(uid uint64, index uint32) string {
	return fmt.Sprintf("m/44'/194'/%d'/0'/%d'", uid, index)
}
