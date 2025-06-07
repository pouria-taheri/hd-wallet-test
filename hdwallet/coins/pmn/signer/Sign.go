package signer

import (
	"encoding/json"
	"fmt"
	ad "git.mazdax.tech/blockchain/hdwallet/account/domain"
	"git.mazdax.tech/blockchain/hdwallet/coins/pmn/domain"
	baseDomain "git.mazdax.tech/blockchain/hdwallet/manager/domain"
	"git.mazdax.tech/data-layer/configcore"
	"git.mazdax.tech/data-layer/loggercore/logger"
	"github.com/aliworkshop/stellar-go/keypair"
	"github.com/aliworkshop/stellar-go/txnbuild"
)

type Signer struct {
	wallet      domain.KuknosWallet
	logger      logger.Logger
	NetworkType string
}

type SignMsg struct {
	Transaction *txnbuild.Transaction `json:"transaction"`
	Account     ad.Account            `json:"account"`
}

type signerConf struct {
	ChainType string
}

func (s Signer) Coin() string {
	return "pmn"
}

func (s Signer) SignTransaction(request []byte) ([]byte, error) {

	var signMsg SignMsg
	err := json.Unmarshal(request, &signMsg)
	if err != nil {
		s.logger.With(logger.Field{
			"submodule": "sign transaction",
			"error":     err,
		}).ErrorF("error in unmarshal transaction request")
		return nil, err
	}

	tx := signMsg.Transaction

	dPath := pathFromUserID(signMsg.Account.Id)
	prv, err := s.wallet.AccPrv(dPath)
	if err != nil {
		s.logger.With(logger.Field{
			"submodule": "sign transaction",
			"error":     err,
			"dPath":     dPath,
		}).ErrorF("error in get account private key")
		return nil, err
	}

	keyPair, err := keypair.FromRawSeed(prv.RawSeed())
	if err != nil {
		s.logger.With(logger.Field{
			"submodule": "sign transaction",
			"error":     err,
			"dPath":     dPath,
		}).ErrorF("error in get account key  pair")
		return nil, err
	}

	s.logger.With(logger.Field{
		"submodule": "sign transaction",
		"payload":   signMsg,
	}).InfoF("signing transaction...")

	signedTx, err := tx.Sign(getNetwork(s.NetworkType), keyPair)
	if err != nil {
		s.logger.With(logger.Field{
			"submodule": "sign transaction",
			"error":     err,
		}).ErrorF("error on signing transaction")
		return nil, err
	}

	res, err := json.Marshal(signedTx)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func NewSigner(logger logger.Logger, conf configcore.Registry, wallet domain.KuknosWallet) baseDomain.SignerModel {
	sc := new(signerConf)
	_ = conf.Unmarshal(sc)
	s := Signer{
		wallet,
		logger,
		sc.ChainType,
	}
	return s
}

func pathFromUserID(uid uint64) string {
	return fmt.Sprintf("m/44'/1'/%d'", uid)
}

func getNetwork(chainType string) string {
	switch chainType {
	case "testnet":
		return "Kuknos-NET"
	case "mainnet":
		return "Kuknos Foundation, Feb 2019"
	}
	return "Kuknos-NET"
}
