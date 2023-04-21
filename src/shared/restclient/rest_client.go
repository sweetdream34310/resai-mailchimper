package restclient

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"net/http"
	"time"

	"gopkg.in/resty.v1"
)

// Client ...
type Client interface {
	SetAddress(address string)
	SetTimeout(timeout time.Duration)
	DefaultHeader(username, password string) http.Header
	BasicAuth(username, password string) string
	Post(path string, headers http.Header, payload interface{}) (body []byte, statusCode int, err error)
	PostWithContext(ctx context.Context, path string, headers http.Header, payload interface{}) (body []byte, statusCode int, err error)
	PostFormData(path string, headers http.Header, payload map[string]string) (body []byte, statusCode int, err error)
	PostFormDataWithContext(ctx context.Context, path string, headers http.Header, payload map[string]string) (body []byte, statusCode int, err error)
	Put(path string, headers http.Header, payload interface{}) (body []byte, statusCode int, err error)
	PutWithContext(ctx context.Context, path string, headers http.Header, payload interface{}) (body []byte, statusCode int, err error)
	Get(path string, headers http.Header) (body []byte, statusCode int, err error)
	GetWithContext(ctx context.Context, path string, headers http.Header) (body []byte, statusCode int, err error)
	GetWithQueryParam(path string, headers http.Header, queryParam map[string]string) (body []byte, statusCode int, err error)
	GetWithQueryParamAndContext(ctx context.Context, path string, headers http.Header, queryParam map[string]string) (body []byte, statusCode int, err error)
	Delete(path string, headers http.Header, payload interface{}) (body []byte, statusCode int, err error)
	DeleteWithContext(ctx context.Context, path string, headers http.Header, payload interface{}) (body []byte, statusCode int, err error)
	Patch(path string, headers http.Header, payload interface{}) (body []byte, statusCode int, err error)
	PatchWithContext(ctx context.Context, path string, headers http.Header, payload interface{}) (body []byte, statusCode int, err error)
}

// New ...
func New(options Options) Client {
	httpClient := resty.New()

	if options.SkipTLS {
		httpClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}

	if options.SkipCheckRedirect {
		httpClient.SetRedirectPolicy(resty.RedirectPolicyFunc(func(request *http.Request, requests []*http.Request) error {
			return http.ErrUseLastResponse
		}))
	}

	if options.WithProxy {
		httpClient.SetProxy(options.ProxyAddress)
	} else {
		httpClient.RemoveProxy()
	}

	httpClient.SetTimeout(options.Timeout * time.Second)
	httpClient.SetDebug(options.DebugMode)

	return &client{
		options:    options,
		httpClient: httpClient,
	}
}

// client ...
type client struct {
	options    Options
	httpClient *resty.Client
}

// DefaultHeader ...
func (c *client) DefaultHeader(username, password string) http.Header {
	headers := http.Header{}
	headers.Set("Authorization", "Basic "+c.BasicAuth(username, password))
	return headers
}

// SetAddress ...
func (c *client) SetAddress(address string) {
	c.options.Address = address
}

// SetTimeout ...
func (c *client) SetTimeout(timeout time.Duration) {
	c.httpClient.SetTimeout(timeout)
}

// BasicAuth ...
func (c *client) BasicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

// Post ...
func (c *client) Post(path string, headers http.Header, payload interface{}) (body []byte, statusCode int, err error) {
	url := c.options.Address + path

	request := c.httpClient.R()
	request.SetBody(payload)

	for h, val := range headers {
		request.Header[h] = val
	}
	if headers["Content-Type"] == nil {
		request.Header.Set("Content-Type", "application/json")
	}

	httpResp, httpErr := request.Post(url)

	if httpResp != nil {
		body = httpResp.Body()
	}

	if httpResp != nil && httpResp.StatusCode() != 0 {
		statusCode = httpResp.StatusCode()
	}

	if statusCode == http.StatusOK {
		return body, statusCode, nil
	}

	return body, statusCode, httpErr
}

// PostWithContext ...
func (c *client) PostWithContext(ctx context.Context, path string, headers http.Header, payload interface{}) (body []byte, statusCode int, err error) {
	url := c.options.Address + path

	request := c.httpClient.R()
	request.SetContext(ctx)
	request.SetBody(payload)

	for h, val := range headers {
		request.Header[h] = val
	}
	if headers["Content-Type"] == nil {
		request.Header.Set("Content-Type", "application/json")
	}

	httpResp, httpErr := request.Post(url)

	if httpResp != nil {
		body = httpResp.Body()
	}

	if httpResp != nil && httpResp.StatusCode() != 0 {
		statusCode = httpResp.StatusCode()
	}

	if statusCode == http.StatusOK {
		return body, statusCode, nil
	}

	return body, statusCode, httpErr
}

// PostFormData ...
func (c *client) PostFormData(path string, headers http.Header, payload map[string]string) (body []byte, statusCode int, err error) {
	url := c.options.Address + path

	request := c.httpClient.R()
	request.SetFormData(payload)

	for h, val := range headers {
		request.Header[h] = val
	}
	if headers["Content-Type"] == nil {
		request.Header.Set("Content-Type", "application/json")
	}

	httpResp, httpErr := request.Post(url)

	if httpResp != nil {
		body = httpResp.Body()
	}

	if httpResp != nil && httpResp.StatusCode() != 0 {
		statusCode = httpResp.StatusCode()
	}

	if statusCode == http.StatusOK {
		return body, statusCode, nil
	}

	return body, statusCode, httpErr
}

// PostFormDataWithContext ...
func (c *client) PostFormDataWithContext(ctx context.Context, path string, headers http.Header, payload map[string]string) (body []byte, statusCode int, err error) {
	url := c.options.Address + path

	request := c.httpClient.R()
	request.SetContext(ctx)
	request.SetFormData(payload)

	for h, val := range headers {
		request.Header[h] = val
	}
	if headers["Content-Type"] == nil {
		request.Header.Set("Content-Type", "application/json")
	}

	httpResp, httpErr := request.Post(url)

	if httpResp != nil {
		body = httpResp.Body()
	}

	if httpResp != nil && httpResp.StatusCode() != 0 {
		statusCode = httpResp.StatusCode()
	}

	if statusCode == http.StatusOK {
		return body, statusCode, nil
	}

	return body, statusCode, httpErr
}

// Put ...
func (c *client) Put(path string, headers http.Header, payload interface{}) (body []byte, statusCode int, err error) {
	url := c.options.Address + path

	request := c.httpClient.R()

	for h, val := range headers {
		request.Header[h] = val
	}
	if headers["Content-Type"] == nil {
		request.Header.Set("Content-Type", "application/json")
	}

	request.SetBody(payload)

	httpResp, httpErr := request.Put(url)

	if httpResp != nil {
		body = httpResp.Body()
	}

	if httpResp != nil && httpResp.StatusCode() != 0 {
		statusCode = httpResp.StatusCode()
	}

	if statusCode == http.StatusOK {
		return body, statusCode, nil
	}

	return body, statusCode, httpErr
}

// PutWithContext ...
func (c *client) PutWithContext(ctx context.Context, path string, headers http.Header, payload interface{}) (body []byte, statusCode int, err error) {
	url := c.options.Address + path

	request := c.httpClient.R()
	request.SetContext(ctx)

	for h, val := range headers {
		request.Header[h] = val
	}
	if headers["Content-Type"] == nil {
		request.Header.Set("Content-Type", "application/json")
	}

	request.SetBody(payload)

	httpResp, httpErr := request.Put(url)

	if httpResp != nil {
		body = httpResp.Body()
	}

	if httpResp != nil && httpResp.StatusCode() != 0 {
		statusCode = httpResp.StatusCode()
	}

	if statusCode == http.StatusOK {
		return body, statusCode, nil
	}

	return body, statusCode, httpErr
}

// Get ...
func (c *client) Get(path string, headers http.Header) (body []byte, statusCode int, err error) {
	url := c.options.Address + path

	request := c.httpClient.R()

	for h, val := range headers {
		request.Header[h] = val
	}

	httpResp, httpErr := request.Get(url)

	if httpResp != nil {
		body = httpResp.Body()
	}

	if httpResp != nil && httpResp.StatusCode() != 0 {
		statusCode = httpResp.StatusCode()
	}

	if statusCode == http.StatusOK {
		return body, statusCode, nil
	}

	return body, statusCode, httpErr
}

// GetWithContext ...
func (c *client) GetWithContext(ctx context.Context, path string, headers http.Header) (body []byte, statusCode int, err error) {
	url := c.options.Address + path

	request := c.httpClient.R()
	request.SetContext(ctx)

	for h, val := range headers {
		request.Header[h] = val
	}

	httpResp, httpErr := request.Get(url)

	if httpResp != nil {
		body = httpResp.Body()
	}

	if httpResp != nil && httpResp.StatusCode() != 0 {
		statusCode = httpResp.StatusCode()
	}

	if statusCode == http.StatusOK {
		return body, statusCode, nil
	}

	return body, statusCode, httpErr
}

// GetWithQueryParam ...
func (c *client) GetWithQueryParam(path string, headers http.Header, queryParam map[string]string) (body []byte, statusCode int, err error) {
	url := c.options.Address + path

	request := c.httpClient.R()

	for h, val := range headers {
		request.Header[h] = val
	}
	request.SetQueryParams(queryParam)

	httpResp, httpErr := request.Get(url)

	if httpResp != nil {
		body = httpResp.Body()
	}

	if httpResp != nil && httpResp.StatusCode() != 0 {
		statusCode = httpResp.StatusCode()
	}

	if statusCode == http.StatusOK {
		return body, statusCode, nil
	}

	return body, statusCode, httpErr
}

// GetWithQueryParamAndContext ...
func (c *client) GetWithQueryParamAndContext(ctx context.Context, path string, headers http.Header, queryParam map[string]string) (body []byte, statusCode int, err error) {
	url := c.options.Address + path

	request := c.httpClient.R()
	request.SetContext(ctx)

	for h, val := range headers {
		request.Header[h] = val
	}
	request.SetQueryParams(queryParam)

	httpResp, httpErr := request.Get(url)

	if httpResp != nil {
		body = httpResp.Body()
	}

	if httpResp != nil && httpResp.StatusCode() != 0 {
		statusCode = httpResp.StatusCode()
	}

	if statusCode == http.StatusOK {
		return body, statusCode, nil
	}

	return body, statusCode, httpErr
}

// Delete ...
func (c *client) Delete(path string, headers http.Header, payload interface{}) (body []byte, statusCode int, err error) {
	url := c.options.Address + path

	request := c.httpClient.R()

	for h, val := range headers {
		request.Header[h] = val
	}

	request.SetBody(payload)

	httpResp, httpErr := request.Delete(url)

	if httpResp != nil {
		body = httpResp.Body()
	}

	if httpResp != nil && httpResp.StatusCode() != 0 {
		statusCode = httpResp.StatusCode()
	}

	if statusCode == http.StatusOK {
		return body, statusCode, nil
	}

	return body, statusCode, httpErr
}

// DeleteWithContext ...
func (c *client) DeleteWithContext(ctx context.Context, path string, headers http.Header, payload interface{}) (body []byte, statusCode int, err error) {
	url := c.options.Address + path

	request := c.httpClient.R()
	request.SetContext(ctx)

	for h, val := range headers {
		request.Header[h] = val
	}

	request.SetBody(payload)

	httpResp, httpErr := request.Delete(url)

	if httpResp != nil {
		body = httpResp.Body()
	}

	if httpResp != nil && httpResp.StatusCode() != 0 {
		statusCode = httpResp.StatusCode()
	}

	if statusCode == http.StatusOK {
		return body, statusCode, nil
	}

	return body, statusCode, httpErr
}

// Patch ...
func (c *client) Patch(path string, headers http.Header, payload interface{}) (body []byte, statusCode int, err error) {
	url := c.options.Address + path

	request := c.httpClient.R()
	request.SetBody(payload)

	for h, val := range headers {
		request.Header[h] = val
	}
	if headers["Content-Type"] == nil {
		request.Header.Set("Content-Type", "application/json")
	}

	httpResp, httpErr := request.Patch(url)

	if httpResp != nil {
		body = httpResp.Body()
	}

	if httpResp != nil && httpResp.StatusCode() != 0 {
		statusCode = httpResp.StatusCode()
	}

	if statusCode == http.StatusOK {
		return body, statusCode, nil
	}

	return body, statusCode, httpErr
}

// PatchWithContext ...
func (c *client) PatchWithContext(ctx context.Context, path string, headers http.Header, payload interface{}) (body []byte, statusCode int, err error) {
	url := c.options.Address + path

	request := c.httpClient.R()
	request.SetContext(ctx)
	request.SetBody(payload)

	for h, val := range headers {
		request.Header[h] = val
	}
	if headers["Content-Type"] == nil {
		request.Header.Set("Content-Type", "application/json")
	}

	httpResp, httpErr := request.Patch(url)

	if httpResp != nil {
		body = httpResp.Body()
	}

	if httpResp != nil && httpResp.StatusCode() != 0 {
		statusCode = httpResp.StatusCode()
	}

	if statusCode == http.StatusOK {
		return body, statusCode, nil
	}

	return body, statusCode, httpErr
}
