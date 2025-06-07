package signer

import (
	"encoding/json"
	ad "git.mazdax.tech/blockchain/hdwallet/account/domain"
	"git.mazdax.tech/blockchain/hdwallet/coins/dot/domain"
	"git.mazdax.tech/blockchain/hdwallet/coins/dot/helper"
	"git.mazdax.tech/blockchain/hdwallet/coins/dot/types"
	baseDomain "git.mazdax.tech/blockchain/hdwallet/manager/domain"
	"git.mazdax.tech/data-layer/configcore"
	"git.mazdax.tech/data-layer/loggercore/logger"
)

type Signer struct {
	wallet    domain.DotWallet
	logger    logger.Logger
	NetworkId types.Network
}

type SignMsg struct {
	Transaction string     `json:"transaction"`
	Account     ad.Account `json:"account"`
}

type signerConf struct {
	NetworkId types.Network
}

func (s Signer) Coin() string {
	return "dot"
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

	dPath := s.dPathFromUserID(signMsg.Account.Id)
	prvKey, err := s.wallet.AccPrv(dPath)
	if err != nil {
		s.logger.With(logger.Field{
			"submodule": "sign transaction",
			"error":     err,
			"dPath":     dPath.String(),
		}).ErrorF("error in get key pair from dPath")
		return nil, err
	}

	txBytes, err := hexDecodeString(signMsg.Transaction)
	if err != nil {
		return nil, err
	}

	s.logger.With(logger.Field{
		"submodule": "sign transaction",
		"payload":   signMsg,
	}).InfoF("signing transaction...")

	signedTx, errSign := prvKey.Sign(txBytes)
	if errSign != nil {
		return nil, errSign
	}

	return u8aConcat([][]byte{{0}, signedTx}), nil
}

func NewSigner(logger logger.Logger, conf configcore.Registry, wallet domain.DotWallet) baseDomain.SignerModel {
	sc := new(signerConf)
	_ = conf.Unmarshal(sc)
	s := Signer{
		wallet,
		logger,
		sc.NetworkId,
	}
	return s
}

func (s Signer) dPathFromUserID(uid uint64) helper.DerivationPath {
	return helper.DerivationPath{
		Network: s.NetworkId,
		Account: uid,
	}
}
