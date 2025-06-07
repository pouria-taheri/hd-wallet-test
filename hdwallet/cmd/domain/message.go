package domain

import (
	"fmt"
	hdwallet "git.mazdax.tech/blockchain/hdwallet/cmd/grpc/domain"
)

type MessageTypeEnum int32

const (
	MessageTypeEnumUnknown  MessageTypeEnum = -1
	MessageTypeEnumGenerate MessageTypeEnum = 0
	MessageTypeEnumEncrypt  MessageTypeEnum = 1
	MessageTypeEnumDecrypt  MessageTypeEnum = 2
	MessageTypeEnumUnlock   MessageTypeEnum = 3

	MessageTypeEnumError         MessageTypeEnum = 100
	MessageTypeEnumMessage       MessageTypeEnum = 101
	MessageTypeEnumTextInput     MessageTypeEnum = 102
	MessageTypeEnumPasswordInput MessageTypeEnum = 103
	MessageTypeEnumFile          MessageTypeEnum = 104

	MessageTypeEnumPing MessageTypeEnum = 999
	MessageTypeEnumPong MessageTypeEnum = 1000
	MessageTypeEnumDone MessageTypeEnum = 200
)

type MessageModel interface {
	SetType(msgType MessageTypeEnum)
	SetArgs(args []string)
	SetData(data []byte)
	SetError(err error)

	GetType() MessageTypeEnum
	GetArgs() []string
	GetMessage() string
	GetData() []byte

	Input(m MessageModel)
	ReadInput() MessageModel
	Output(m MessageModel)
	ReadOutput() MessageModel
	Close()

	GetGrpcBody() *hdwallet.Body
}

type Message struct {
	Type    MessageTypeEnum
	Args    []string
	Message string
	Data    []byte
	Err     error

	inputChan  chan MessageModel
	outputChan chan MessageModel
}

func NewMessage() *Message {
	return &Message{
		inputChan:  make(chan MessageModel, 1),
		outputChan: make(chan MessageModel, 1),
	}
}

func NewMessageModel() MessageModel {
	return NewMessage()
}

func (msg *Message) Input(m MessageModel) {
	msg.inputChan <- m
}

func (msg *Message) ReadInput() MessageModel {
	m := <-msg.inputChan
	return m
}

func (msg *Message) Output(m MessageModel) {
	msg.outputChan <- m
}

func (msg *Message) ReadOutput() MessageModel {
	m := <-msg.outputChan
	return m
}

func (msg *Message) Close() {
	close(msg.outputChan)
	close(msg.inputChan)
}

func (msg *Message) SetType(msgType MessageTypeEnum) {
	msg.Type = msgType
}

func (msg *Message) SetArgs(args []string) {
	msg.Args = args
}

func (msg *Message) SetData(data []byte) {
	msg.Data = data
}

func (msg *Message) SetError(err error) {
	msg.Err = err
	msg.Message = err.Error()
}

func (msg *Message) GetType() MessageTypeEnum {
	return msg.Type
}

func (msg *Message) GetArgs() []string {
	return msg.Args
}

func (msg *Message) GetMessage() string {
	return msg.Message
}

func (msg *Message) GetData() []byte {
	return msg.Data
}

func (msg *Message) GetGrpcBody() *hdwallet.Body {
	body := &hdwallet.Body{
		Type:    0,
		Message: msg.Message,
		Data:    msg.Data,
		Args:    msg.Args,
	}
	switch msg.Type {
	case MessageTypeEnumGenerate:
		body.Type = hdwallet.MessageTypeEnum_Generate
	case MessageTypeEnumEncrypt:
		body.Type = hdwallet.MessageTypeEnum_Encrypt
	case MessageTypeEnumDecrypt:
		body.Type = hdwallet.MessageTypeEnum_Decrypt
	case MessageTypeEnumUnlock:
		body.Type = hdwallet.MessageTypeEnum_Unlock
	case MessageTypeEnumError:
		body.Type = hdwallet.MessageTypeEnum_Error
	case MessageTypeEnumMessage:
		body.Type = hdwallet.MessageTypeEnum_Message
	case MessageTypeEnumTextInput:
		body.Type = hdwallet.MessageTypeEnum_TextInput
	case MessageTypeEnumPasswordInput:
		body.Type = hdwallet.MessageTypeEnum_PasswordInput
	case MessageTypeEnumFile:
		body.Type = hdwallet.MessageTypeEnum_File
	case MessageTypeEnumDone:
		body.Type = hdwallet.MessageTypeEnum_Done
	case MessageTypeEnumPing:
		body.Type = hdwallet.MessageTypeEnum_Ping
	case MessageTypeEnumPong:
		body.Type = hdwallet.MessageTypeEnum_Pong
	}
	return body
}

func Sprintf(str string, a ...interface{}) MessageModel {
	msg := NewMessage()
	msg.Type = MessageTypeEnumMessage
	msg.Message = fmt.Sprintf(str+"\n", a...)
	return msg
}

func SprintFln(str string, a ...interface{}) MessageModel {
	msg := NewMessage()
	msg.Type = MessageTypeEnumMessage
	msg.Message = fmt.Sprintf(str+"\n", a...)
	return msg
}

func ReadPassword(str string, a ...interface{}) MessageModel {
	msg := NewMessage()
	msg.Type = MessageTypeEnumPasswordInput
	msg.Message = fmt.Sprintf(str, a...)
	return msg
}

func Error(err error) *Message {
	msg := NewMessage()
	msg.Type = MessageTypeEnumError
	msg.Message = err.Error() + "\n"
	msg.Err = err
	return msg
}

func GetMessageFromHdWallet(body *hdwallet.Body) MessageModel {
	msg := NewMessage()
	msg.Args = body.Args
	msg.Message = body.Message
	msg.Data = body.Data

	switch body.Type {
	case hdwallet.MessageTypeEnum_Unknown:
		msg.Type = MessageTypeEnumUnknown
	case hdwallet.MessageTypeEnum_Generate:
		msg.Type = MessageTypeEnumGenerate
	case hdwallet.MessageTypeEnum_Encrypt:
		msg.Type = MessageTypeEnumEncrypt
	case hdwallet.MessageTypeEnum_Decrypt:
		msg.Type = MessageTypeEnumDecrypt
	case hdwallet.MessageTypeEnum_Unlock:
		msg.Type = MessageTypeEnumUnlock
	case hdwallet.MessageTypeEnum_Error:
		msg.Type = MessageTypeEnumError
	case hdwallet.MessageTypeEnum_Message:
		msg.Type = MessageTypeEnumMessage
	case hdwallet.MessageTypeEnum_TextInput:
		msg.Type = MessageTypeEnumTextInput
	case hdwallet.MessageTypeEnum_PasswordInput:
		msg.Type = MessageTypeEnumPasswordInput
	case hdwallet.MessageTypeEnum_File:
		msg.Type = MessageTypeEnumFile
	case hdwallet.MessageTypeEnum_Done:
		msg.Type = MessageTypeEnumDone
	case hdwallet.MessageTypeEnum_Ping:
		msg.Type = MessageTypeEnumPing
	case hdwallet.MessageTypeEnum_Pong:
		msg.Type = MessageTypeEnumPong
	}

	return msg
}
