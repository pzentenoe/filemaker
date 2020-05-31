package filemaker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type RecordService interface {
	Create(payload *Payload) (*ResponseData, error)
	Edit(recordId string, payload *Payload) (*ResponseData, error)
	Duplicate(recordId string) (*ResponseData, error)
	Delete(recordId string) (*ResponseData, error)
	GetById(recordId string) (*ResponseData, error)
	List(offset, limit string, sorters ...*Sorter) (*ResponseData, error)
}

const recordsPath = "fmi/data/%s/databases/%s/layouts/%s/records"

type recordService struct {
	database string
	layout   string
	client   *Client
}

func NewRecordService(database, layout string, client *Client) *recordService {
	return &recordService{
		database: database,
		layout:   layout,
		client:   client,
	}
}

type Payload struct {
	FieldData  interface{} `json:"fieldData"`
	PortalData interface{} `json:"portalData,omitempty"`
}

func (s *recordService) Create(payload *Payload) (*ResponseData, error) {

	responseAuth, err := s.client.Connect(s.database)
	if err != nil {
		return nil, err
	}

	defer s.client.Disconnect(s.database, responseAuth.Response.Token)

	path := fmt.Sprintf(recordsPath, s.client.version, s.database, s.layout)
	options := &performRequestOptions{
		Method:      http.MethodPost,
		Path:        path,
		ContentType: "application/json",
		Body:        payload,
		Headers: http.Header{
			"Authorization": []string{fmt.Sprintf("Bearer %s", responseAuth.Response.Token)},
		},
	}

	return s.client.executeQuery(options)
}

func (s *recordService) Edit(recordId string, payload *Payload) (*ResponseData, error) {

	responseAuth, err := s.client.Connect(s.database)
	if err != nil {
		return nil, err
	}

	defer s.client.Disconnect(s.database, responseAuth.Response.Token)

	path := fmt.Sprintf(recordsPath+"/%s", s.client.version, s.database, s.layout, recordId)
	options := &performRequestOptions{
		Method:      http.MethodPatch,
		Path:        path,
		ContentType: "application/json",
		Body:        payload,
		Headers: http.Header{
			"Authorization": []string{fmt.Sprintf("Bearer %s", responseAuth.Response.Token)},
		},
	}

	return s.client.executeQuery(options)

}

func (s *recordService) Duplicate(recordId string) (*ResponseData, error) {
	responseAuth, err := s.client.Connect(s.database)
	if err != nil {
		return nil, err
	}

	defer s.client.Disconnect(s.database, responseAuth.Response.Token)

	path := fmt.Sprintf(recordsPath+"/%s", s.client.version, s.database, s.layout, recordId)
	options := &performRequestOptions{
		Method:      http.MethodPost,
		Path:        path,
		ContentType: "application/json",
		Headers: http.Header{
			"Authorization": []string{fmt.Sprintf("Bearer %s", responseAuth.Response.Token)},
		},
	}

	return s.client.executeQuery(options)
}

func (s *recordService) Delete(recordId string) (*ResponseData, error) {

	responseAuth, err := s.client.Connect(s.database)
	if err != nil {
		return nil, err
	}

	defer s.client.Disconnect(s.database, responseAuth.Response.Token)

	path := fmt.Sprintf(recordsPath+"/%s", s.client.version, s.database, s.layout, recordId)
	options := &performRequestOptions{
		Method: http.MethodDelete,
		Path:   path,
		Headers: http.Header{
			"Authorization": []string{fmt.Sprintf("Bearer %s", responseAuth.Response.Token)},
		},
	}

	return s.client.executeQuery(options)

}
func (s *recordService) GetById(recordId string) (*ResponseData, error) {

	responseAuth, err := s.client.Connect(s.database)
	if err != nil {
		return nil, err
	}

	defer s.client.Disconnect(s.database, responseAuth.Response.Token)

	path := fmt.Sprintf(recordsPath+"/%s", s.client.version, s.database, s.layout, recordId)
	options := &performRequestOptions{
		Method: http.MethodGet,
		Path:   path,
		Headers: http.Header{
			"Authorization": []string{fmt.Sprintf("Bearer %s", responseAuth.Response.Token)},
		},
	}

	return s.client.executeQuery(options)

}

func (s *recordService) List(offset, limit string, sorters ...*Sorter) (*ResponseData, error) {
	responseAuth, err := s.client.Connect(s.database)
	if err != nil {
		return nil, err
	}

	defer s.client.Disconnect(s.database, responseAuth.Response.Token)

	path := fmt.Sprintf(recordsPath, s.client.version, s.database, s.layout)

	params := url.Values{}
	params.Add("_offset", offset)
	params.Add("_limit", limit)

	sortersStr := sortersToJson(sorters...)
	if sortersStr != "" {
		params.Add("_sort", sortersStr)
	}

	options := &performRequestOptions{
		Method: http.MethodGet,
		Path:   path,
		Params: params,
		Headers: http.Header{
			"Authorization": []string{fmt.Sprintf("Bearer %s", responseAuth.Response.Token)},
		},
	}
	return s.client.executeQuery(options)
}

func sortersToJson(sorters ...*Sorter) string {
	if len(sorters) > 0 {
		data, err := json.Marshal(sorters)
		if err != nil {
			return ""
		}
		return string(data)
	}
	return ""
}
