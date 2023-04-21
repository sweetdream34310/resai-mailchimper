package cbc

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

func Encrypt(key, plainText []byte, opts ...Opts) (cipherText []byte, err error) {
	pad := PKCS5Padding
	if len(opts) > 0 {
		pad = opts[0].Pad
	}
	if len(plainText)%aes.BlockSize != 0 {
		plainText = pad(plainText, aes.BlockSize)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	cipherText = make([]byte, aes.BlockSize+len(plainText))

	iv := cipherText[:aes.BlockSize]

	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return
	}

	cbc := cipher.NewCBCEncrypter(block, iv)
	cbc.CryptBlocks(cipherText[aes.BlockSize:], plainText)

	return
}

func Decrypt(key, cipherText []byte, opts ...Opts) (plainText []byte, err error) {
	var block cipher.Block

	unpad := PKCS5Trimming
	if len(opts) > 0 {
		unpad = opts[0].Unpad
	}

	if block, err = aes.NewCipher(key); err != nil {
		return
	}

	iv := cipherText[:aes.BlockSize]

	cipherText = cipherText[aes.BlockSize:]

	cbc := cipher.NewCBCDecrypter(block, iv)
	cbc.CryptBlocks(cipherText, cipherText)

	plainText = unpad(cipherText)

	return
}
