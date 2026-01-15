package filemaker

import (
	"fmt"
	"net/http"
)

type SearchService interface {
	Do() (interface{}, error)
}

const findQueryPath = "fmi/data/%s/databases/%s/layouts/%s/_find"

type searchService struct {
	client    *Client
	database  string
	layout    string
	seachData *searchData
}

func NewSearchService(database, layout string, client *Client) *searchService {
	searchData := new(searchData)
	searchData.QueryGroup = make([]map[string]string, 0)
	return &searchService{
		client:    client,
		database:  database,
		layout:    layout,
		seachData: searchData,
	}
}

type searchData struct {
	QueryGroup []map[string]string `json:"query"`
	Limit      string              `json:"limit,omitempty"`
	Offset     string              `json:"offset,omitempty"`
	Portal     []string            `json:"portal"`
	Sort       []*Sorter           `json:"sort,omitempty"`
}

func (s *searchService) GroupQueries(queryGroups ...*groupQuery) *searchService {
	queries := make([]map[string]string, 0)
	for _, queryGroup := range queryGroups {
		queryMap := make(map[string]string)
		for _, query := range queryGroup.queries {
			value := query.valueWithOp()
			queryMap[query.Name] = value
		}
		queries = append(queries, queryMap)
	}
	s.seachData.QueryGroup = queries
	return s
}

func (s *searchService) SetOffset(offset string) *searchService {
	s.seachData.Offset = offset
	return s
}

func (s *searchService) SetLimit(limit string) *searchService {
	s.seachData.Limit = limit
	return s
}

func (s *searchService) SetPortals(portal []string) *searchService {
	s.seachData.Portal = portal
	return s
}

func (s *searchService) Sorters(sorters ...*Sorter) *searchService {
	s.seachData.Sort = sorters
	return s
}

func (s *searchService) Do() (*ResponseData, error) {

	responseAuth, err := s.client.Connect(s.database)
	if err != nil {
		return nil, err
	}
	defer s.client.Disconnect(s.database, responseAuth.Response.Token)

	path := fmt.Sprintf(findQueryPath, s.client.version, s.database, s.layout)

	options := &performRequestOptions{
		Method:  http.MethodPost,
		Path:    path,
		Body:    s.seachData,
		Headers: http.Header{"Authorization": []string{fmt.Sprintf("Bearer %s", responseAuth.Response.Token)}},
	}

	return s.client.executeQuery(options)
}
