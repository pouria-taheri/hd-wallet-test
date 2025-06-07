package wallet

import (
	"fmt"
	"git.mazdax.tech/blockchain/hdwallet/coins/ftm/domain"
	"git.mazdax.tech/data-layer/loggercore/logger"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"strconv"
)

func createWalletMainAccount(key *hdkeychain.ExtendedKey, log logger.Logger) domain.WalletMainAcc {
	wma := domain.WalletMainAcc{}

	dPath, err := pathFromUserID(0, 0)
	if err != nil {
		log.With(logger.Field{
			"submodule": "walletMainAccount",
			"userId":    0,
			"index":     0,
			"error":     err,
		}).ErrorF("error in get path from userId and index")
		panic("error in get path")
	}
	wma.Path = dPath
	log.With(logger.Field{
		"submodule": "walletMainAccount",
		"dPath":     dPath,
	}).InfoF("Main account")

	prv, err := privateFromMasterAndAddressBTC(key, wma.Path)
	if err != nil {
		log.With(logger.Field{
			"submodule": "walletMainAccount",
			"error":     err,
			"path":      wma.Path,
		}).ErrorF("error in get private key")
		panic("error in get private key from key and path")
	}

	wma.PrivateKey = prv.ToECDSA()

	addr := addressFromPrivate(prv)
	wma.Address = &addr
	log.With(logger.Field{
		"submodule": "walletMainAccount",
		"address":   wma.Address.String(),
	}).InfoF("Main account")

	wma.PublicKey = prv.PubKey().ToECDSA()

	return wma
}

func privateFromMasterAndAddressBTC(key *hdkeychain.ExtendedKey, path accounts.DerivationPath) (*btcec.PrivateKey, error) {
	var err error
	for _, n := range path {
		key, err = key.Child(n)
		if err != nil {
			return nil, err
		}
	}
	return key.ECPrivKey()
}

func pathFromUserID(uid uint64, index uint32) (accounts.DerivationPath, error) {
	str := strconv.FormatUint(uid, 10)
	path := fmt.Sprintf("m/44/1007/%s/0/%d", str, index)
	dPath, err := accounts.ParseDerivationPath(path)
	if err != nil {
		return accounts.DerivationPath{}, err
	}
	return dPath, nil
}

func addressFromPrivate(key *btcec.PrivateKey) common.Address {
	pub := key.PubKey().ToECDSA()
	return crypto.PubkeyToAddress(*pub)
}
