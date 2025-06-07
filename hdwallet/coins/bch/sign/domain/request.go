package domain

import (
	bd "git.mazdax.tech/blockchain/hdwallet/coins/bch/domain"
	"github.com/bchsuite/bchd/wire"
	"github.com/bchsuite/bchutil"
)

type SignRequest struct {
	Tx               *wire.MsgTx            `json:"tx"`
	InputValues      []bchutil.Amount       `json:"inputValues"`
	InputAddresses   []*bd.Address          `json:"inputAddresses"`
	ExtractedScripts []bd.ExtractedPKScript `json:"extractedScripts"`
	PrevScripts      [][]byte               `json:"prevScripts"`
}
