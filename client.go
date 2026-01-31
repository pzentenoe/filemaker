package filemaker

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"runtime"
	"strings"
	"sync"
)

type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

type performRequestOptions struct {
	Method      string
	Path        string
	Params      url.Values
	Body        any
	ContentType string
	Headers     http.Header
	Username    string
	Password    string
	BasicAuth   bool
}

type Client struct {
	mu           sync.RWMutex
	url          string //URL with port or dns
	version      string //Default vLatest
	httpClient   *http.Client
	retryConfig  *RetryConfig
	logger       Logger
	metrics      *Metrics
	authProvider AuthProvider
	username     string
	password     string
}

// getVersion safely retrieves the client version with read lock.
func (c *Client) getVersion() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.version
}

// getURL safely retrieves the client URL with read lock.
func (c *Client) getURL() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.url
}

func NewClient(options ...ClientOptions) (*Client, error) {
	c := &Client{}
	for _, option := range options {
		if err := option(c); err != nil {
			return nil, err
		}
	}
	if c.version == "" {
		c.version = DefaultVersion
	}
	if c.httpClient == nil {
		c.httpClient = http.DefaultClient
	}
	if c.retryConfig == nil {
		// Use default retry configuration
		c.retryConfig = DefaultRetryConfig()
	}
	if c.logger == nil {
		// Use no-op logger by default
		c.logger = &NoOpLogger{}
	}
	// Metrics are optional and remain nil unless explicitly enabled

	return c, nil
}

func (c *Client) executeQuery(ctx context.Context, options *performRequestOptions) (*ResponseData, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	var responseData *ResponseData

	err := c.retryConfig.executeWithRetry(ctx, func() error {
		response, err := c.performRequest(ctx, options)

		if response == nil && err != nil {
			return fmt.Errorf("failed to perform request: %w", err)
		}

		if response != nil {
			defer func() { _ = response.Body.Close() }()

			var parseErr error
			responseData, parseErr = c.parseResponse(response)
			return parseErr
		}

		return nil
	})

	return responseData, err
}

func (c *Client) parseResponse(response *http.Response) (*ResponseData, error) {
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, &NetworkError{
			Message: "failed to read response body",
			Err:     err,
		}
	}

	var responseData ResponseData
	err = json.Unmarshal(data, &responseData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(responseData.Messages) > 0 {
		if fmErr := ParseFileMakerError(&responseData, response.StatusCode); fmErr != nil {
			return &responseData, fmErr
		}
	}

	if response.StatusCode >= 400 {
		return &responseData, &FileMakerError{
			HTTPStatus: response.StatusCode,
			Message:    http.StatusText(response.StatusCode),
		}
	}

	return &responseData, nil
}

func (c *Client) performRequest(ctx context.Context, opt *performRequestOptions) (*http.Response, error) {

	if c.url == "" {
		return nil, &ValidationError{
			Field:   "url",
			Message: "URL is required",
		}
	}
	if opt.Method == "" {
		return nil, &ValidationError{
			Field:   "method",
			Message: "HTTP method is required",
		}
	}

	pathWithParams := opt.Path
	if len(opt.Params) > 0 {
		pathWithParams += "?" + opt.Params.Encode()
	}

	completeUrl := fmt.Sprintf("%s/%s", c.url, pathWithParams)

	req, err := c.NewRequest(opt.Method, completeUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	if opt.ContentType != "" {
		req.Header.Set("Content-Type", opt.ContentType)
	}
	if len(opt.Headers) > 0 {
		for key, value := range opt.Headers {
			for _, v := range value {
				req.Header.Add(key, v)
			}
		}
	}

	if opt.Body != nil {
		err = req.setBody(opt.Body, false)
		if err != nil {
			return nil, fmt.Errorf("failed to set request body: %w", err)
		}
	}

	if opt.BasicAuth {
		req.setBasicAuth(opt.Username, opt.Password)
	}

	resp, err := c.Do((*http.Request)(req).WithContext(ctx))
	return resp, err

}

type Request http.Request

func (c *Client) NewRequest(method, url string) (*Request, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", "filemaker/"+c.version+" ("+runtime.GOOS+"-"+runtime.GOARCH+")")
	req.Header.Add("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	return (*Request)(req), nil
}

func (r *Request) setBody(body any, gzipCompress bool) error {
	switch b := body.(type) {
	case string:
		if gzipCompress {
			return r.setBodyGzip(b)
		}
		return r.setBodyString(b)
	default:
		if gzipCompress {
			return r.setBodyGzip(body)
		}
		return r.setBodyJson(body)
	}
}

func (r *Request) setBodyJson(data any) error {
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}
	r.Header.Set("Content-Type", "application/json")
	return r.setBodyReader(bytes.NewReader(body))
}

func (r *Request) setBodyGzip(body any) error {
	switch b := body.(type) {
	case string:
		buf := new(bytes.Buffer)
		w := gzip.NewWriter(buf)
		if _, err := w.Write([]byte(b)); err != nil {
			return err
		}
		if err := w.Close(); err != nil {
			return err
		}
		r.Header.Add("Content-Encoding", "gzip")
		r.Header.Add("Vary", "Accept-Encoding")
		return r.setBodyReader(bytes.NewReader(buf.Bytes()))
	default:
		data, err := json.Marshal(b)
		if err != nil {
			return err
		}
		buf := new(bytes.Buffer)
		w := gzip.NewWriter(buf)
		if _, err := w.Write(data); err != nil {
			return err
		}
		if err := w.Close(); err != nil {
			return err
		}
		r.Header.Add("Content-Encoding", "gzip")
		r.Header.Add("Vary", "Accept-Encoding")
		r.Header.Set("Content-Type", "application/json")
		return r.setBodyReader(bytes.NewReader(buf.Bytes()))
	}
}

func (r *Request) setBodyString(body string) error {
	return r.setBodyReader(strings.NewReader(body))
}

func (r *Request) setBodyReader(body io.Reader) error {
	rc, ok := body.(io.ReadCloser)
	if !ok && body != nil {
		rc = io.NopCloser(body)
	}
	r.Body = rc
	if body != nil {
		switch v := body.(type) {
		case *strings.Reader:
			r.ContentLength = int64(v.Len())
		case *bytes.Buffer:
			r.ContentLength = int64(v.Len())
		}
	}
	return nil
}

func (r *Request) setBasicAuth(username, password string) {
	((*http.Request)(r)).SetBasicAuth(username, password)
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	return c.httpClient.Do(req)
}
