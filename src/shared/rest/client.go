package rest

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

// RequestOptions - Request Headers
type RequestOptions struct {
	Payload interface{}
	Headers map[string]string
}

type Response struct {
	Error       bool        `json:"error"`
	Data        interface{} `json:"data"`
	Status      int         `json:"status"`
	UserMessage string      `json:"userMessage"`
	Errors      interface{} `json:"errors"`
}

// Request - Request
func Request(method, requestURL string, options *RequestOptions) (response Response, err error) {
	var request *http.Request
	client := &http.Client{}
	if method != "GET" {
		reqBody, err := json.Marshal(options.Payload)
		if err != nil {
			return response, fmt.Errorf("could not marshal request or invalid json: %v", err)
		}
		request, err = http.NewRequest(method, requestURL, bytes.NewBuffer(reqBody))
	} else {
		if options.Payload != nil {
			payload := &bytes.Buffer{}
			writer := multipart.NewWriter(payload)
			for k, v := range options.Payload.(map[string]string) {
				_ = writer.WriteField(k, v)
			}
			err := writer.Close()
			if err != nil {
				fmt.Println(err)
			}
			request, err = http.NewRequest(method, requestURL, payload)
			request.Header.Set("Content-Type", writer.FormDataContentType())
		} else {
			request, err = http.NewRequest(method, requestURL, nil)
		}
	}
	if err != nil {
		return response, fmt.Errorf("could not create request: %v", err)
	}
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept-Encoding", "gzip")
	if options.Headers != nil {
		for key, value := range options.Headers {
			request.Header.Add(key, value)
		}
	}
	resp, err := client.Do(request)
	if err != nil {
		return response, fmt.Errorf("unable to make request: %v", err)
	}

	reader := resp.Body
	if resp.Header.Get("Content-Encoding") == "gzip" {
		reader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return response, fmt.Errorf("could not read gzip body from response: %v", err)
		}
		defer reader.Close()
	}

	resBody, err := ioutil.ReadAll(reader)
	if err != nil {
		return response, fmt.Errorf("could not read body from response: %v", err)
	}

	if err = json.Unmarshal(resBody, &response); err != nil {
		return response, fmt.Errorf("unable to get the correct response: %v", err)
	}

	response.Status = resp.StatusCode

	return response, nil
}
