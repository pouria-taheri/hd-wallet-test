package domain

type Session struct {
	Input  chan MessageModel
	Output chan MessageModel
}
