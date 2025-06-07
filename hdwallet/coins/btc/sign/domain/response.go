package domain

import "github.com/btcsuite/btcd/wire"

type SignResponse struct {
	Tx *wire.MsgTx `json:"tx"`
}
