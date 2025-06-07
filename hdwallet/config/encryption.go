package config

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
)

type UnlockHandler func(encrypted []byte) (unencrypted []byte, err error)

var (
	lockedKey = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
)

func GetValidKey(salt string, password []byte) []byte {
	saltBytes := []byte(salt)
	// ensure key length is 32 chars(aes-256)
	uselessKeyLen := len(salt) + len(password) - 32
	if uselessKeyLen < 0 {
		uselessKeyLen = -uselessKeyLen
		saltBytes = append(saltBytes, []byte(strings.Repeat(salt,
			uselessKeyLen/len(salt)+1))...)
		uselessKeyLen = len(saltBytes) + len(password) - 32
	}
	saltBytes = saltBytes[:len(saltBytes)-uselessKeyLen]
	key := append(saltBytes, password...)
	return key
}

func encodeBase64(b []byte) []byte {
	return []byte(base64.StdEncoding.EncodeToString(b))
}

func decodeBase64(b []byte) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(string(b))
	if err != nil {
		return nil, err
	}
	return data, nil
}

func Encrypt(data, key []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	b := encodeBase64(data)
	ciphertext := make([]byte, aes.BlockSize+len(b))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], b)
	return ciphertext
}

func Decrypt(data, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(data) < aes.BlockSize {
		err := errors.New(fmt.Sprintf("length of data must be at least "+
			"%d bytes", aes.BlockSize))
		return nil, err
	}
	iv := data[:aes.BlockSize]
	data = data[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(data, data)
	result, err := decodeBase64(data)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func EncryptFile(source, destination string, key []byte) {
	data, err := ioutil.ReadFile(source)
	if err != nil {
		panic(err)
	}
	ciphertext := Encrypt(data, key)
	// insert key at first of file to determine if file is encrypted or not
	toSave := append(lockedKey, ciphertext...)
	err = ioutil.WriteFile(destination, toSave, 0644)
	if err != nil {
		panic(err)
	}
}

func DecryptFile(path string, key []byte) (io.Reader, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if len(data) < 10 {
		return nil, errors.New("no enough data in file")
	}
	data = data[10:]
	result, err := Decrypt(data, key)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(result), err
}
