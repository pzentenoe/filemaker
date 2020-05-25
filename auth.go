package filemaker

import (
	"fmt"
	"net/http"
)

const (
	sessionAuthPath string = "fmi/data/%s/databases/%s/sessions"
)

type ConnectionDatasource struct {
	FmDataSource []FmDatasource `json:"fmDataSource"`
}
type FmDatasource struct {
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (c *Client) Connect(database string) (*ResponseData, error) {
	c.mu.RLock()
	path := fmt.Sprintf(sessionAuthPath, c.version, database)

	options := &performRequestOptions{
		Method:      http.MethodPost,
		Path:        path,
		Body:        "{}",
		ContentType: "application/json",
		basicAuth:   true,
	}

	response, err := c.executeQuery(options)
	c.mu.RUnlock()
	return response, err
}

func (c *Client) ConnectWithDatasource(database string) (*ResponseData, error) {
	c.mu.RLock()
	path := fmt.Sprintf(sessionAuthPath, c.version, database)
	datasource := FmDatasource{
		Database: database,
		Username: c.username,
		Password: c.password,
	}

	fileMakerConnection := ConnectionDatasource{FmDataSource: []FmDatasource{
		datasource,
	}}

	options := &performRequestOptions{
		Method:    http.MethodPost,
		Path:      path,
		Body:      fileMakerConnection,
		basicAuth: true,
	}
	response, err := c.executeQuery(options)
	c.mu.RUnlock()
	return response, err
}

func (c *Client) Disconnect(database, token string) (*ResponseData, error) {
	c.mu.RLock()
	path := fmt.Sprintf(sessionAuthPath+"/%s", c.version, database, token)

	options := &performRequestOptions{
		Method: http.MethodDelete,
		Path:   path,
	}
	response, err := c.executeQuery(options)
	c.mu.RUnlock()
	return response, err
}
