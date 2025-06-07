package domain

type HandlerModel interface {
	Generate(msg MessageModel)
	Encrypt(msg MessageModel)
	Decrypt(msg MessageModel)
	Unlock(msg MessageModel)
	LoadSecureConfigs()
	EnsureSecureConfigLoaded()
	WaitForUnlock()
}
