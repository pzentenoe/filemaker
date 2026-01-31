package filemaker

import (
	"context"
	"fmt"
	"net/http"
)

// SearchService defines the interface for performing find operations in FileMaker.
type SearchService interface {
	Do(ctx context.Context) (*ResponseData, error)
}

const findQueryPath = "fmi/data/%s/databases/%s/layouts/%s/_find"

type searchService struct {
	client     *Client
	database   string
	layout     string
	searchData *searchData
}

// NewSearchService creates a new instance of SearchService for a specific database and layout.
func NewSearchService(database, layout string, client *Client) *searchService {
	searchData := new(searchData)
	searchData.QueryGroup = make([]map[string]string, 0)
	return &searchService{
		client:     client,
		database:   database,
		layout:     layout,
		searchData: searchData,
	}
}

type searchData struct {
	QueryGroup []map[string]string `json:"query"`
	Limit      string              `json:"limit,omitempty"`
	Offset     string              `json:"offset,omitempty"`
	Portal     []string            `json:"portal"`
	Sort       []*Sorter           `json:"sort,omitempty"`
}

// GroupQueries adds query groups to the search criteria.
// Multiple query groups are joined with OR logic.
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
	s.searchData.QueryGroup = queries
	return s
}

// SetOffset sets the starting record number for pagination.
func (s *searchService) SetOffset(offset string) *searchService {
	s.searchData.Offset = offset
	return s
}

// SetLimit sets the maximum number of records to return.
func (s *searchService) SetLimit(limit string) *searchService {
	s.searchData.Limit = limit
	return s
}

// SetPortals specifies which portal data to include in the response.
func (s *searchService) SetPortals(portal []string) *searchService {
	s.searchData.Portal = portal
	return s
}

// Sorters sets the sort order for the search results.
func (s *searchService) Sorters(sorters ...*Sorter) *searchService {
	s.searchData.Sort = sorters
	return s
}

// Do executes the search query and returns the matching records.
func (s *searchService) Do(ctx context.Context) (*ResponseData, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	responseAuth, err := s.client.CreateSession(ctx, s.database)
	if err != nil {
		return nil, err
	}
	defer func(client *Client, ctx context.Context, database, token string) {
		_, _ = client.DisconnectWithContext(ctx, database, token)

	}(s.client, ctx, s.database, responseAuth.Response.Token)

	path := fmt.Sprintf(findQueryPath, s.client.getVersion(), s.database, s.layout)

	options := &performRequestOptions{
		Method:  http.MethodPost,
		Path:    path,
		Body:    s.searchData,
		Headers: http.Header{"Authorization": []string{fmt.Sprintf("Bearer %s", responseAuth.Response.Token)}},
	}

	return s.client.executeQuery(ctx, options)
}
