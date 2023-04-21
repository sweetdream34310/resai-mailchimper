package google

import (
	"crypto/rand"
	"errors"
	"io"
	"net/http"
	"os"

	"github.com/cloudsrc/api.awaymail.v1.go/src/shared/constants"
)

func (p *provider) getClientIDSecret(agent string) (clientID, clientSecret string, err error) {
	switch agent {
	case constants.IosAgent:
		clientID = p.config.Google.Ios.ClientID
		clientSecret = p.config.Google.Ios.ClientSecret
	case constants.WebAgent:
		clientID = p.config.Google.Website.ClientID
		clientSecret = p.config.Google.Website.ClientSecret
	default:
		err = constants.ErrorInvalidRequest
		return
	}
	return
}

func randStr(strSize int, randType string) string {

	var dictionary string

	if randType == "alphanum" {
		dictionary = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	}

	var strBytes = make([]byte, strSize)
	_, _ = rand.Read(strBytes)
	for k, v := range strBytes {
		strBytes[k] = dictionary[v%byte(len(dictionary))]
	}
	return string(strBytes)
}

func downloadFile(URL, fileName string) error {
	//Get the response bytes from the url
	response, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return errors.New("file or url not found")
	}
	//Create a empty file
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	//Write the bytes to the fiel
	if _, err := io.Copy(file, response.Body); err != nil {
		return err
	}
	return nil
}

func chunkSplit(body string, limit int, end string) string {
	var charSlice []rune

	// push characters to slice
	for _, char := range body {
		charSlice = append(charSlice, char)
	}

	var result = ""

	for len(charSlice) >= 1 {
		// convert slice/array back to string
		// but insert end at specified limit
		result = result + string(charSlice[:limit]) + end

		// discard the elements that were copied over to result
		charSlice = charSlice[limit:]

		// change the limit
		// to cater for the last few words in
		if len(charSlice) < limit {
			limit = len(charSlice)
		}
	}
	return result
}
