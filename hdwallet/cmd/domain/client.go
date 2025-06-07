package domain

type ClientModel interface {
	SetInnerClient(client ClientModel)

	Generate(msg MessageModel)
	Encrypt(msg MessageModel)
	Decrypt(msg MessageModel)
	Unlock(msg MessageModel)

	HandleMessage(msg MessageModel) MessageModel
	Handle(msg MessageModel)
}
