package domain

import (
	bd "git.mazdax.tech/blockchain/hdwallet/coins/btc/domain"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
)

type SignRequest struct {
	Tx               *wire.MsgTx            `json:"tx"`
	InputValues      []btcutil.Amount       `json:"inputValues"`
	InputAddresses   []*bd.Address          `json:"inputAddresses"`
	ExtractedScripts []bd.ExtractedPKScript `json:"extractedScripts"`
	PrevScripts      [][]byte               `json:"prevScripts"`
}
