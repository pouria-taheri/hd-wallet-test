package wallet

import (
	"fmt"
	"git.mazdax.tech/blockchain/hdwallet/coins/sol/domain"
	"git.mazdax.tech/data-layer/loggercore/logger"
	"github.com/gagliardetto/solana-go"
)

func createWalletMainAccount(key domain.MasterKey, log logger.Logger) domain.WalletMainAcc {
	wma := domain.WalletMainAcc{}

	wma.PrivateKey = key.MasterKey.Derive(144).Derive(0).Derive(0)

	wma.PublicKey = solana.PublicKeyFromBytes(wma.PrivateKey.PublicKey())

	log.With(logger.Field{
		"submodule": "walletMainAccount",
		"address":   wma.PublicKey,
	}).InfoF("Main account")

	return wma
}

func pathFromUserId(uid uint64, index uint32) string {
	return fmt.Sprintf("m/44'/144'/%d'/0'/%d'", uid, index)
}
