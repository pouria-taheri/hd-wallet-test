package signer

import (
	"encoding/json"
	"fmt"
	ad "git.mazdax.tech/blockchain/hdwallet/account/domain"
	"git.mazdax.tech/blockchain/hdwallet/coins/avax/domain"
	"git.mazdax.tech/data-layer/configcore"
	"git.mazdax.tech/data-layer/loggercore/logger"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
	"strconv"
)

type Signer struct {
	Wallet    domain.AvaxWallet
	eipSigner types.Signer
	logger    logger.Logger
}

type SignMsg struct {
	Transaction *types.Transaction `json:"transaction"`
	Account     *ad.Account        `json:"account"`
}

type signerConf struct {
	NetworkId int64
}

func (s Signer) Coin() string {
	return "avax"
}

func (s Signer) SignTransaction(request []byte) ([]byte, error) {
	signMsg := new(SignMsg)
	err := json.Unmarshal(request, signMsg)
	if err != nil {
		return nil, err
	}

	tx := signMsg.Transaction
	dPath, err := pathFromUserID(signMsg.Account.Id, signMsg.Account.Index)
	if err != nil {
		return nil, err
	}

	privateKey, err := s.Wallet.AccPrv(dPath)
	if err != nil {
		return nil, err
	}

	signed, err := types.SignTx(tx, s.eipSigner, privateKey)
	if err != nil {
		return nil, err
	}

	return json.Marshal(signed)
}

func NewSigner(logger logger.Logger, conf configcore.Registry, wallet domain.AvaxWallet) Signer {
	sc := new(signerConf)
	_ = conf.Unmarshal(sc)
	s := Signer{
		Wallet:    wallet,
		eipSigner: types.NewEIP155Signer(big.NewInt(sc.NetworkId)),
		logger:    logger,
	}
	return s
}

func pathFromUserID(uid uint64, index uint32) (accounts.DerivationPath, error) {
	str := strconv.FormatUint(uid, 10)
	path := fmt.Sprintf("m/44/9005/%s/0/%d", str, index)
	dPath, err := accounts.ParseDerivationPath(path)
	if err != nil {
		return accounts.DerivationPath{}, err
	}
	return dPath, nil
}
