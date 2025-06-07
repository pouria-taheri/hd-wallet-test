package signer

import (
	"encoding/json"
	"fmt"
	"git.mazdax.tech/blockchain/bnb-go-sdk/common/types"
	gtypes "git.mazdax.tech/blockchain/bnb-go-sdk/types"
	"git.mazdax.tech/blockchain/bnb-go-sdk/types/msg"
	"git.mazdax.tech/blockchain/bnb-go-sdk/types/tx"
	ad "git.mazdax.tech/blockchain/hdwallet/account/domain"
	"git.mazdax.tech/blockchain/hdwallet/coins/bnb/domain"
	baseDomain "git.mazdax.tech/blockchain/hdwallet/manager/domain"
	"git.mazdax.tech/data-layer/configcore"
	"git.mazdax.tech/data-layer/loggercore/logger"
)

type Transaction struct {
	Msg           msg.SendMsg
	AccountNumber int64 `json:"account_number"`
	Sequence      int64 `json:"sequence"`
}

type Signer struct {
	wallet    domain.BinanceWallet
	NetworkId types.ChainNetwork
	logger    logger.Logger
}

type SignMsg struct {
	Transaction *Transaction `json:"transaction"`
	Account     ad.Account   `json:"account"`
}

type signerConf struct {
	NetworkId types.ChainNetwork
}

func (s *Signer) Coin() string {
	return "bnb"
}

func (s *Signer) SignTransaction(request []byte) ([]byte, error) {

	signMsg := new(SignMsg)
	err := json.Unmarshal(request, signMsg)
	if err != nil {
		s.logger.With(logger.Field{
			"submodule": "sign transaction",
			"error":     err,
		}).ErrorF("error in unmarshal transaction request")
		return nil, err
	}
	keyManager := s.wallet.AccKeyManager(s.logger, signMsg.Account)
	if keyManager == nil {
		s.logger.With(logger.Field{
			"submodule": "sign transaction",
			"accountId": signMsg.Account.Id,
			"index":     signMsg.Account.Index,
		}).WarnF("error in get account key manager")
		return nil, fmt.Errorf("keymanager is missing")
	}

	transaction := signMsg.Transaction

	// prepare message to sign
	chainID := gtypes.ProdChainID
	if s.NetworkId == types.TestNetwork {
		chainID = gtypes.TestnetChainID
	} else if s.NetworkId == types.TmpTestNetwork {
		chainID = gtypes.KongoChainId
	} else if s.NetworkId == types.GangesNetwork {
		chainID = gtypes.GangesChainId
	}
	stdSignMsg := &tx.StdSignMsg{
		ChainID:       chainID,
		AccountNumber: transaction.AccountNumber,
		Sequence:      transaction.Sequence,
		Memo:          signMsg.Account.Memo,
		Msgs:          []msg.Msg{transaction.Msg},
		Source:        tx.Source,
	}
	s.logger.With(logger.Field{
		"submodule": "sign transaction",
		"payload":   signMsg,
	}).InfoF("signing transaction...")

	//for _, op := range transaction.options {
	//	signMsg = op(signMsg)
	//}

	for _, m := range stdSignMsg.Msgs {
		if err := m.ValidateBasic(); err != nil {
			s.logger.With(logger.Field{
				"submodule": "sign transaction",
				"msg":       m,
			}).WarnF("error in validate transaction message")
			return nil, err
		}
	}

	return keyManager.Sign(*stdSignMsg)
}

func NewSigner(logger logger.Logger, conf configcore.Registry, wallet domain.BinanceWallet) baseDomain.SignerModel {
	sc := new(signerConf)
	_ = conf.Unmarshal(sc)
	s := &Signer{
		wallet:    wallet,
		NetworkId: sc.NetworkId,
		logger:    logger,
	}
	return s
}
