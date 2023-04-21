package cbc

import (
	"bytes"
)

type PadFn func([]byte, int) []byte
type TrimFn func([]byte) []byte

type Opts struct {
	Pad   PadFn
	Unpad TrimFn
}

func PKCS5Padding(chioherText []byte, blockSize int) []byte {
	padding := blockSize - len(chioherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(chioherText, padText...)
}

func PKCS5Trimming(encrypt []byte) []byte {
	padding := encrypt[len(encrypt)-1]
	if len(encrypt)-int(padding) < 1 {
		return encrypt[:]
	}
	return encrypt[:len(encrypt)-int(padding)]
}
