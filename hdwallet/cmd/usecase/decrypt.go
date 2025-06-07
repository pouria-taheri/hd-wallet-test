package usecase

import (
	"git.mazdax.tech/blockchain/hdwallet/cmd/domain"
	"git.mazdax.tech/blockchain/hdwallet/config"
	"github.com/fatih/structs"
	"time"
)

func (h *handler) getPassword(msg domain.MessageModel) []byte {
	var pw []byte
	for {
		msg.Output(domain.ReadPassword("password: "))
		req := msg.ReadInput()
		args := req.GetArgs()
		if len(args) == 0 {
			continue
		}
		pw = []byte(args[0])
		if len(pw) > 0 {
			break
		}
	}
	return pw
}

func (h *handler) decryptDataWithPassword(msg domain.MessageModel, data []byte, pw []byte) ([]byte, error) {
	for {
		key := config.GetValidKey(h.app.Config.Salt, pw)
		var encrypted = make([]byte, len(data))
		copy(encrypted, data)
		decrypted, err := config.Decrypt(encrypted, key)
		if err != nil {
			msg.Output(domain.Error(err))
			return nil, err
		}
		return decrypted, nil
	}
}

func (h *handler) decryptDataWithRequest(msg domain.MessageModel, data []byte) ([]byte, error) {
	for {
		var pw []byte
		pw = h.getPassword(msg)

		return h.decryptDataWithPassword(msg, data, pw)
	}
}

func (h *handler) decryptData(data []byte) ([]byte, error) {
	for {
		unlockRequest := <-h.unlockRequestsChan
		if unlockRequest == nil {
			continue
		}
		return h.decryptDataWithRequest(unlockRequest, data)
	}
}

func (h *handler) getSingleInputDecryptDataHandler(request domain.MessageModel) func(data []byte) ([]byte, error) {
	pwErrored := false
	pw := h.getPassword(request)
	return func(data []byte) ([]byte, error) {
		if pwErrored {
			pw = h.getPassword(request)
		}
		r, err := h.decryptDataWithPassword(request, data, pw)
		if err != nil {
			pwErrored = true
			return nil, err
		}
		return r, nil
	}
}

func (h *handler) getDecryptDataHandler(request domain.MessageModel) func(data []byte) ([]byte, error) {
	return func(data []byte) ([]byte, error) {
		return h.decryptDataWithRequest(request, data)
	}
}

func (h *handler) Decrypt(msg domain.MessageModel) {
	if !config.IsAnyConfigSecured(h.app.SecuredConfigFileDetails) {
		h.app.Logger.InfoF("configuration is already decrypted")
		return
	}

	go h.client.Handle(msg)

	var err error
	defer func() {
		if err != nil {
			msg.Output(domain.Error(err))
		}
	}()
	h.app.SecureConfigs, err = config.LoadSecureConfigs(h.app.SecuredConfigFileDetails, h.getSingleInputDecryptDataHandler(msg))
	if err != nil {
		return
	}

	encryptedPaths := make(map[string]*domain.Message)
	for _, sc := range h.app.SecureConfigs {
		detail := sc.GetDetail()
		if _, ok := encryptedPaths[detail.FilePath]; ok {
			continue
		}
		h.saveConfig(structs.Map(sc), detail.FilePath)
		dataBytes, err := h.loadConfig(detail.FilePath)
		if err != nil {
			return
		}
		encryptedPaths[detail.FilePath] = &domain.Message{
			Type: domain.MessageTypeEnumFile,
			Data: domain.NewFile(detail.Coin, dataBytes),
		}
	}
	msg.Output(domain.SprintFln("done"))
	time.Sleep(time.Second)
	msg.Close()
}
