package wallet

import (
	"crypto/ecdsa"
	"fmt"
	"git.mazdax.tech/blockchain/hdwallet/coins/trx/domain"
	"git.mazdax.tech/data-layer/loggercore/logger"
	"github.com/btcsuite/btcd/btcec"
	"github.com/fbsobreira/gotron-sdk/pkg/address"
	"github.com/fbsobreira/gotron-sdk/pkg/keys/hd"
	"github.com/fbsobreira/gotron-sdk/pkg/keystore"
	"github.com/pborman/uuid"
	"strconv"
)

func createWalletMainAccount(key domain.MasterKey, log logger.Logger) domain.WalletMainAcc {
	wma := domain.WalletMainAcc{}

	dPath := pathFromUserID(0, 0)
	wma.Path = dPath
	log.With(logger.Field{
		"submodule": "walletMainAccount",
		"dPath":     dPath,
	}).InfoF("Main account")

	prv, pub, err := privateFromMasterAndAddress(key.Secret, key.ChainCode, wma.Path)
	if err != nil {
		log.With(logger.Field{
			"submodule": "walletMainAccount",
			"error":     err,
			"path":      wma.Path,
		}).ErrorF("error in get private key")
		panic("error in get private key from secret and chainCode and path")
	}

	wma.PrivateKey = prv.ToECDSA()
	wma.PublicKey = pub.ToECDSA()

	addr := addressFromPrivate(prv)
	wma.Account = &addr
	log.With(logger.Field{
		"submodule": "walletMainAccount",
		"address":   wma.Account.Address.String(),
	}).InfoF("Main account")

	return wma
}

func privateFromMasterAndAddress(master, ch [32]byte, path string) (*btcec.PrivateKey, *btcec.PublicKey, error) {
	private, err := hd.DerivePrivateKeyForPath(
		btcec.S256(),
		master,
		ch,
		path,
	)
	if err != nil {
		return nil, nil, err
	}
	sk, pk := btcec.PrivKeyFromBytes(btcec.S256(), private[:])
	return sk, pk, nil
}

func pathFromUserID(uid uint64, index uint32) string {
	str := strconv.FormatUint(uid, 10)
	return fmt.Sprintf("44'/195'/%s'/0/%d", str, index)
}

func addressFromPrivate(privateKeyECDSA *btcec.PrivateKey) keystore.Account {
	id := uuid.NewRandom()
	key := &keystore.Key{
		ID:         id,
		Address:    address.PubkeyToAddress(privateKeyECDSA.PublicKey),
		PrivateKey: privateKeyECDSA.ToECDSA(),
	}
	return keystore.Account{Address: key.Address, URL: keystore.URL{Scheme: keystore.KeyStoreScheme}}
}

func (w tronWallet) MasterAccPrv() *ecdsa.PrivateKey { return w.WalletMainAcc.PrivateKey }
