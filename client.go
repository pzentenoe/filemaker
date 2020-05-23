package filemaker

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
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
	Body        interface{}
	ContentType string
	Headers     http.Header
	basicAuth   bool
}

type Client struct {
	mu         sync.RWMutex
	url        string //URL with port or dns
	username   string
	password   string
	version    string //Default vLatest
	httpClient *http.Client
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

	return c, nil
}

func (c *Client) executeQuery(options *performRequestOptions) (*ResponseData, error) {
	response, err := c.performRequest(context.Background(), options)

	var searchResponseData *ResponseData
	if response == nil && err != nil {
		fmt.Errorf("Fail performRequest %s", err.Error())
		return nil, err
	} else if response != nil && err != nil {
		defer response.Body.Close()

		errFillResponse := fillSearchResponseFromHttpResponse(response, searchResponseData)
		if errFillResponse != nil {
			return nil, err
		}
	} else if err == nil {
		defer response.Body.Close()

		err := fillSearchResponseFromHttpResponse(response, searchResponseData)
		if err != nil {
			return nil, err
		}
	}

	return searchResponseData, nil
}

func fillSearchResponseFromHttpResponse(response *http.Response, searchResponseData *ResponseData) error {
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Errorf("Error to get response body : %s", err.Error())
		return err
	}

	err = json.Unmarshal(data, &searchResponseData)
	if err != nil {
		fmt.Errorf("Error to Unmarshal QueryResponse %s", err.Error())
		return err
	}
	return nil
}

func (c *Client) performRequest(ctx context.Context, opt *performRequestOptions) (*http.Response, error) {

	if c.url == "" {
		return nil, errors.New("Empty URL")
	}
	if opt.Method == "" {
		return nil, errors.New("Empty Method")
	}

	pathWithParams := opt.Path
	if len(opt.Params) > 0 {
		pathWithParams += "?" + opt.Params.Encode()
	}

	completeUrl := fmt.Sprintf("%s/%s", c.url, pathWithParams)

	req, err := c.NewRequest(opt.Method, completeUrl)
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
			return nil, fmt.Errorf("filemaker: couldn't set body %v for request: %v", opt.Body, err)
		}
	}

	if opt.basicAuth {
		req.setBasicAuth(c.username, c.password)
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

func (r *Request) setBody(body interface{}, gzipCompress bool) error {
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

func (r *Request) setBodyJson(data interface{}) error {
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}
	r.Header.Set("Content-Type", "application/json")
	r.setBodyReader(bytes.NewReader(body))
	return nil
}

func (r *Request) setBodyGzip(body interface{}) error {
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
		rc = ioutil.NopCloser(body)
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
