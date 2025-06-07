package domain

import "encoding/json"

type File struct {
	Name string
	Data []byte
}

func NewFile(name string, data []byte) []byte {
	f := File{
		Name: name,
		Data: data,
	}
	bytes, _ := json.Marshal(f)
	return bytes
}
