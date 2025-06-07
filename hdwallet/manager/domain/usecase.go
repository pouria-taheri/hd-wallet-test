package domain

type UseCaseModel interface {
	RegisterSigner(signer SignerModel)
	RegisterAccountHandler(signer SignerModel)
}
