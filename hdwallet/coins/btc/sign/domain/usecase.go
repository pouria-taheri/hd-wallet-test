package domain

import "git.mazdax.tech/core/errors"

type UseCaseModel interface {
	Coin() string
	SignTransaction(request []byte) ([]byte, error)
	Sign(request SignRequest) (*SignResponse, errors.ErrorModel)
}
