package utils

import (
	"crypto/md5"
	"encoding/hex"

	"github.com/cloudsrc/api.awaymail.v1.go/src/shared/utils/cbc"
)

func Encrypt(salt, plainText string) (cipherText string, err error) {
	hasher := md5.New()
	hasher.Write([]byte(salt))
	key := hex.EncodeToString(hasher.Sum(nil))

	cipher, err := cbc.Encrypt([]byte(key), []byte(plainText))
	if err != nil {
		return
	}
	cipherText = hex.EncodeToString(cipher)
	return
}

func Decrypt(salt, cipherText string) (plainText string, err error) {
	hasher := md5.New()
	hasher.Write([]byte(salt))
	key := hex.EncodeToString(hasher.Sum(nil))

	cipher, err := hex.DecodeString(cipherText)
	if err != nil {
		return
	}

	plain, err := cbc.Decrypt([]byte(key), cipher)
	if err != nil {
		return
	}
	plainText = string(plain)

	return
}
