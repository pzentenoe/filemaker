package filemaker

import (
	"errors"
	"net/http"
)

const (
	DefaultVersion = "vLatest"
)

type ClientOptions func(*client) error

func SetURL(url string) ClientOptions {
	return func(c *client) error {
		if url == "" {
			return errors.New("Empty url")
		}
		c.url = url
		return nil
	}
}

func SetUsername(username string) ClientOptions {
	return func(c *client) error {
		if username == "" {
			return errors.New("Empty username")
		}
		c.username = username
		return nil
	}
}

func SetPassword(password string) ClientOptions {
	return func(c *client) error {
		if password == "" {
			return errors.New("Empty password")
		}
		c.password = password
		return nil
	}
}

func SetVersion(version string) ClientOptions {
	return func(c *client) error {
		if version == "" {
			c.version = DefaultVersion
		}
		c.version = version
		return nil
	}
}

func SetHttpClient(httpClient *http.Client) ClientOptions {
	return func(c *client) error {
		if httpClient == nil {
			c.httpClient = http.DefaultClient
		}
		c.httpClient = httpClient
		return nil
	}
}