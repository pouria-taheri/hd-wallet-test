package signer

import (
	"encoding/json"
	ad "git.mazdax.tech/blockchain/hdwallet/account/domain"
	"git.mazdax.tech/blockchain/hdwallet/coins/atom/domain"
	baseDomain "git.mazdax.tech/blockchain/hdwallet/manager/domain"
	"git.mazdax.tech/data-layer/configcore"
	"git.mazdax.tech/data-layer/loggercore/logger"
	c "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
)

type Signer struct {
	wallet      domain.CosmosWallet
	logger      logger.Logger
	NetworkType string
}

type Transaction struct {
	From             string `json:"from"`
	To               string `json:"to"`
	Value            int64  `json:"value"`
	Memo             string `json:"memo"`
	SignaturePayload string `json:"signature"`
}

type SigningPayload struct {
	Address       string        `json:"address"`
	Bytes         []byte        `json:"bytes"`
	SignatureType SignatureType `json:"signature_type,omitempty"`
}

type SignatureType string

type UnsignedTx struct {
	SigningPayload *SigningPayload
	TxBuilder      c.TxBuilder
	TxConfig       c.TxConfig
}

type SignMsg struct {
	SigningPayload []byte     `json:"transaction"`
	Account        ad.Account `json:"account"`
}

type signerConf struct {
	ChainType string
}

func (s Signer) Coin() string {
	return "atom"
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

	path := hd.NewFundraiserParams(uint32(signMsg.Account.Id), 118, 0)
	prv, err := s.wallet.AccPrv(path)
	if err != nil {
		s.logger.With(logger.Field{
			"submodule": "sign transaction",
			"error":     err,
			"dPath":     path,
		}).ErrorF("error in get account private key")
		return nil, err
	}

	s.logger.With(logger.Field{
		"submodule": "sign transaction",
		"payload":   signMsg,
	}).InfoF("signing transaction...")

	sk, err := prv.Sign(signMsg.SigningPayload)
	if err != nil {
		return nil, err
	}

	return sk, nil
}

func NewSigner(logger logger.Logger, conf configcore.Registry, wallet domain.CosmosWallet) baseDomain.SignerModel {
	sc := new(signerConf)
	_ = conf.Unmarshal(sc)
	s := Signer{
		wallet,
		logger,
		sc.ChainType,
	}
	return s
}
