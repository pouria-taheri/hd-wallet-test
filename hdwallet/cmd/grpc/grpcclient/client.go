package grpcclient

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"git.mazdax.tech/blockchain/hdwallet/cmd/domain"
	hdwallet "git.mazdax.tech/blockchain/hdwallet/cmd/grpc/domain"
	"git.mazdax.tech/blockchain/hdwallet/config"
	"git.mazdax.tech/data-layer/loggercore/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"time"
)

type Model interface {
	domain.ClientModel
	Initialize()
}

type client struct {
	config config.Cmd
	logger logger.Logger
	client domain.ClientModel

	conn *grpc.ClientConn
}

func New(config config.Cmd, logger logger.Logger) Model {
	c := &client{
		config: config,
		logger: logger,
	}
	return c
}

func (c *client) loadTLSCredentials() (credentials.TransportCredentials, error) {
	if c.config.Client.Grpc.CAFile == "" {
		return nil, nil
	}
	pemServerCA, err := ioutil.ReadFile(c.config.Client.Grpc.CAFile)
	if err != nil {
		return nil, err
	}
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, fmt.Errorf("failed to add server CA's certificate")
	}
	cfg := &tls.Config{
		RootCAs: certPool,
	}
	return credentials.NewTLS(cfg), nil
}

func (c *client) Initialize() {
	dialOptions := []grpc.DialOption{grpc.WithBlock()}
	cred, err := c.loadTLSCredentials()
	if err != nil {
		panic(err)
	}
	var insecure bool
	if cred != nil {
		dialOptions = append(dialOptions, grpc.WithTransportCredentials(cred))
	} else {
		insecure = true
	}
	if insecure {
		dialOptions = append(dialOptions, grpc.WithInsecure())
	}

	for {
		ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
		c.logger.InfoF("connecting to server to address %s...",
			c.config.Client.Grpc.Address)
		var err error
		c.conn, err = grpc.DialContext(ctx, c.config.Client.Grpc.Address, dialOptions...)
		if err == nil {
			c.logger.InfoF("connected.")
			break
		}
		c.logger.InfoF("connection could not be established, retrying...")
		time.Sleep(time.Second)
	}
}

func (c *client) SetInnerClient(client domain.ClientModel) {
	c.client = client
}

func (c *client) Generate(msg domain.MessageModel) {
	panic("not implemented")
}

func (c *client) Encrypt(msg domain.MessageModel) {
	panic("not implemented")
}

func (c *client) Decrypt(msg domain.MessageModel) {
	panic("not implemented")
}

func (c *client) innerClientHandleMessage(msg domain.MessageModel) domain.MessageModel {
	if c.client != nil {
		return c.client.HandleMessage(msg)
	}
	return msg
}

func (c *client) HandleMessage(msg domain.MessageModel) domain.MessageModel {
	var result domain.MessageModel
	switch msg.GetType() {
	case domain.MessageTypeEnumUnlock:
		c.Unlock(msg)
		return nil
	}
	result = c.innerClientHandleMessage(msg)
	if result != nil {
		r2 := c.HandleMessage(result)
		if r2 != nil {
			return r2
		}
	}
	return result
}

func (c *client) Handle(msg domain.MessageModel) {
	if c.client != nil {
		c.client.Handle(msg)
	}
}

func (c *client) getMsgTypeFromHdWallet(msgType hdwallet.MessageTypeEnum) domain.MessageTypeEnum {
	switch msgType {
	case hdwallet.MessageTypeEnum_Generate:
		return domain.MessageTypeEnumGenerate
	case hdwallet.MessageTypeEnum_Encrypt:
		return domain.MessageTypeEnumEncrypt
	case hdwallet.MessageTypeEnum_Decrypt:
		return domain.MessageTypeEnumDecrypt
	case hdwallet.MessageTypeEnum_Unlock:
		return domain.MessageTypeEnumUnlock
	case hdwallet.MessageTypeEnum_Error:
		return domain.MessageTypeEnumError
	case hdwallet.MessageTypeEnum_TextInput:
		return domain.MessageTypeEnumTextInput
	case hdwallet.MessageTypeEnum_PasswordInput:
		return domain.MessageTypeEnumPasswordInput
	case hdwallet.MessageTypeEnum_File:
		return domain.MessageTypeEnumFile
	case hdwallet.MessageTypeEnum_Done:
		return domain.MessageTypeEnumDone
	}
	return domain.MessageTypeEnumUnknown
}

func (c *client) getHdWalletMsgTypeFromRequest(msgType domain.MessageTypeEnum) hdwallet.MessageTypeEnum {
	switch msgType {
	case domain.MessageTypeEnumGenerate:
		return hdwallet.MessageTypeEnum_Generate
	case domain.MessageTypeEnumEncrypt:
		return hdwallet.MessageTypeEnum_Encrypt
	case domain.MessageTypeEnumDecrypt:
		return hdwallet.MessageTypeEnum_Decrypt
	case domain.MessageTypeEnumUnlock:
		return hdwallet.MessageTypeEnum_Unlock
	case domain.MessageTypeEnumError:
		return hdwallet.MessageTypeEnum_Error
	case domain.MessageTypeEnumTextInput:
		return hdwallet.MessageTypeEnum_TextInput
	case domain.MessageTypeEnumPasswordInput:
		return hdwallet.MessageTypeEnum_PasswordInput
	case domain.MessageTypeEnumFile:
		return hdwallet.MessageTypeEnum_File
	case domain.MessageTypeEnumDone:
		return hdwallet.MessageTypeEnum_Done
	}
	return hdwallet.MessageTypeEnum_Unknown
}
