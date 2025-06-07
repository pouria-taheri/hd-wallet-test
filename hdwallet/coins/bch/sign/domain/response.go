package domain

import "github.com/bchsuite/bchd/wire"

type SignResponse struct {
	Tx *wire.MsgTx `json:"tx"`
}
