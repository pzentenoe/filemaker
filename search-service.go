package filemaker

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// SearchService defines the interface for performing find operations in FileMaker.
type SearchService interface {
	Do(ctx context.Context) (*ResponseData, error)
}

const findQueryPath = "fmi/data/%s/databases/%s/layouts/%s/_find"

// PortalConfig defines pagination settings for a specific portal.
type PortalConfig struct {
	Name   string // Portal name
	Offset int    // Starting record number (1-based)
	Limit  int    // Maximum number of records to return
}

// NewPortalConfig creates a new PortalConfig with the specified name.
func NewPortalConfig(name string) *PortalConfig {
	return &PortalConfig{
		Name:   name,
		Offset: 1,
		Limit:  50, // Default limit
	}
}

// WithOffset sets the offset for the portal.
func (pc *PortalConfig) WithOffset(offset int) *PortalConfig {
	pc.Offset = offset
	return pc
}

// WithLimit sets the limit for the portal.
func (pc *PortalConfig) WithLimit(limit int) *PortalConfig {
	pc.Limit = limit
	return pc
}

type searchService struct {
	client        *Client
	database      string
	layout        string
	searchData    *searchData
	portalConfigs []*PortalConfig
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
// For portal pagination, use SetPortalConfigs instead.
func (s *searchService) SetPortals(portal []string) *searchService {
	s.searchData.Portal = portal
	return s
}

// SetPortalConfigs configures portals with pagination settings.
// This allows you to specify offset and limit for each portal individually.
//
// Example:
//
//	service.SetPortalConfigs(
//	    NewPortalConfig("RelatedOrders").WithOffset(1).WithLimit(10),
//	    NewPortalConfig("RelatedPayments").WithOffset(1).WithLimit(5),
//	)
func (s *searchService) SetPortalConfigs(configs ...*PortalConfig) *searchService {
	s.portalConfigs = configs

	// Extract portal names for the portal field in JSON body
	portalNames := make([]string, len(configs))
	for i, config := range configs {
		portalNames[i] = config.Name
	}
	s.searchData.Portal = portalNames

	return s
}

// Sorters sets the sort order for the search results.
func (s *searchService) Sorters(sorters ...*Sorter) *searchService {
	s.searchData.Sort = sorters
	return s
}

// Do executes the search query and returns the matching records.
func (s *searchService) Do(ctx context.Context) (*ResponseData, error) {
	ctx = ensureContext(ctx)

	responseAuth, err := s.client.CreateSession(ctx, s.database)
	if err != nil {
		return nil, err
	}
	defer func(client *Client, ctx context.Context, database, token string) {
		_, _ = client.DisconnectWithContext(ctx, database, token)

	}(s.client, ctx, s.database, responseAuth.Response.Token)

	path := fmt.Sprintf(findQueryPath, s.client.getVersion(), s.database, s.layout)

	// Build query parameters for portal pagination
	params := url.Values{}
	for _, portalConfig := range s.portalConfigs {
		if portalConfig.Offset > 0 {
			params.Set("_offset."+portalConfig.Name, strconv.Itoa(portalConfig.Offset))
		}
		if portalConfig.Limit > 0 {
			params.Set("_limit."+portalConfig.Name, strconv.Itoa(portalConfig.Limit))
		}
	}

	options := &performRequestOptions{
		Method:  http.MethodPost,
		Path:    path,
		Params:  params,
		Body:    s.searchData,
		Headers: bearerAuthHeader(responseAuth.Response.Token),
	}

	return s.client.executeQuery(ctx, options)
}
