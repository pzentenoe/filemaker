package filemaker

import (
	"context"
	"strconv"
)

// RecordBuilder provides a fluent interface for building and executing record operations.
type RecordBuilder struct {
	client        *Client
	database      string
	layout        string
	recordID      string
	fieldData     map[string]any
	portalData    map[string][]map[string]any
	scripts       *ScriptContext
	ctx           context.Context
	modID         string
	deleteRelated string
	limit         string
	offset        string
}

// NewRecordBuilder creates a new RecordBuilder for the specified database and layout.
// Example:
//
//	builder := filemaker.NewRecordBuilder(client, "Contacts", "ContactDetails")
//	response, err := builder.
//	    SetField("FirstName", "John").
//	    SetField("LastName", "Doe").
//	    SetField("Email", "john@example.com").
//	    Create(ctx)
func NewRecordBuilder(client *Client, database, layout string) *RecordBuilder {
	return &RecordBuilder{
		client:     client,
		database:   database,
		layout:     layout,
		fieldData:  make(map[string]any),
		portalData: make(map[string][]map[string]any),
		ctx:        context.Background(),
	}
}

// WithContext sets the context for the operation.
func (rb *RecordBuilder) WithContext(ctx context.Context) *RecordBuilder {
	rb.ctx = ctx
	return rb
}

// Limit sets the maximum number of records to return for List operations.
func (rb *RecordBuilder) Limit(limit int) *RecordBuilder {
	rb.limit = strconv.Itoa(limit)
	return rb
}

// Offset sets the starting record for List operations.
func (rb *RecordBuilder) Offset(offset int) *RecordBuilder {
	rb.offset = strconv.Itoa(offset)
	return rb
}

// SetField sets a field value in the record.
func (rb *RecordBuilder) SetField(name string, value any) *RecordBuilder {
	rb.fieldData[name] = value
	return rb
}

// SetFields sets multiple field values at once.
func (rb *RecordBuilder) SetFields(fields map[string]any) *RecordBuilder {
	for name, value := range fields {
		rb.fieldData[name] = value
	}
	return rb
}

// AddPortalRecord adds a portal record to the specified portal.
// Portal records are related records that appear in portals on the layout.
func (rb *RecordBuilder) AddPortalRecord(portalName string, record map[string]any) *RecordBuilder {
	if rb.portalData[portalName] == nil {
		rb.portalData[portalName] = make([]map[string]any, 0)
	}
	rb.portalData[portalName] = append(rb.portalData[portalName], record)
	return rb
}

// SetPortalRecords sets all portal records for the specified portal.
func (rb *RecordBuilder) SetPortalRecords(portalName string, records []map[string]any) *RecordBuilder {
	rb.portalData[portalName] = records
	return rb
}

// WithScripts configures scripts to execute with the operation.
func (rb *RecordBuilder) WithScripts(scripts *ScriptContext) *RecordBuilder {
	rb.scripts = scripts
	return rb
}

// WithPreRequestScript sets a script to execute before the request.
func (rb *RecordBuilder) WithPreRequestScript(script, param string) *RecordBuilder {
	if rb.scripts == nil {
		rb.scripts = NewScriptContext()
	}
	rb.scripts.WithPreRequest(script, param)
	return rb
}

// WithAfterScript sets a script to execute after the request.
func (rb *RecordBuilder) WithAfterScript(script, param string) *RecordBuilder {
	if rb.scripts == nil {
		rb.scripts = NewScriptContext()
	}
	rb.scripts.WithAfter(script, param)
	return rb
}

// ForRecord sets the record ID for update/delete operations.
func (rb *RecordBuilder) ForRecord(recordID string) *RecordBuilder {
	rb.recordID = recordID
	return rb
}

// WithModID sets the modification ID for optimistic locking.
// The operation will only succeed if the record's modID matches.
func (rb *RecordBuilder) WithModID(modID string) *RecordBuilder {
	rb.modID = modID
	return rb
}

// WithDeleteRelated sets whether to delete related records (portal name).
func (rb *RecordBuilder) WithDeleteRelated(deleteRelated string) *RecordBuilder {
	rb.deleteRelated = deleteRelated
	return rb
}

// Create creates a new record with the configured field data.
func (rb *RecordBuilder) Create(ctx context.Context) (*ResponseData, error) {
	if ctx == nil {
		ctx = rb.ctx
	}

	service := NewRecordService(rb.database, rb.layout, rb.client)

	payload := &Payload{
		FieldData: rb.fieldData,
	}

	if len(rb.portalData) > 0 {
		payload.PortalData = rb.portalData
	}

	return service.Create(ctx, payload)
}

// Update updates the record with the configured field data.
// Requires ForRecord() to be called first.
func (rb *RecordBuilder) Update(ctx context.Context) (*ResponseData, error) {
	if ctx == nil {
		ctx = rb.ctx
	}

	if rb.recordID == "" {
		return nil, &ValidationError{
			Field:   "recordID",
			Message: "record ID is required for update operations",
		}
	}

	service := NewRecordService(rb.database, rb.layout, rb.client)

	payload := &Payload{
		FieldData: rb.fieldData,
		ModId:     rb.modID,
	}

	if len(rb.portalData) > 0 {
		payload.PortalData = rb.portalData
	}

	return service.Edit(ctx, rb.recordID, payload)
}

// Delete deletes the record.
// Requires ForRecord() to be called first.
func (rb *RecordBuilder) Delete(ctx context.Context) (*ResponseData, error) {
	if ctx == nil {
		ctx = rb.ctx
	}

	if rb.recordID == "" {
		return nil, &ValidationError{
			Field:   "recordID",
			Message: "record ID is required for delete operations",
		}
	}

	service := NewRecordService(rb.database, rb.layout, rb.client)
	return service.Delete(ctx, rb.recordID, rb.deleteRelated)
}

// Get retrieves the record by ID.
// Requires ForRecord() to be called first.
func (rb *RecordBuilder) Get(ctx context.Context) (*ResponseData, error) {
	if ctx == nil {
		ctx = rb.ctx
	}

	if rb.recordID == "" {
		return nil, &ValidationError{
			Field:   "recordID",
			Message: "record ID is required for get operations",
		}
	}

	service := NewRecordService(rb.database, rb.layout, rb.client)
	return service.GetById(ctx, rb.recordID)
}

// Duplicate duplicates the record.
// Requires ForRecord() to be called first.
func (rb *RecordBuilder) Duplicate(ctx context.Context) (*ResponseData, error) {
	if ctx == nil {
		ctx = rb.ctx
	}

	if rb.recordID == "" {
		return nil, &ValidationError{
			Field:   "recordID",
			Message: "record ID is required for duplicate operations",
		}
	}

	service := NewRecordService(rb.database, rb.layout, rb.client)
	return service.Duplicate(ctx, rb.recordID)
}

// List retrieves a list of records.
// Use Limit() and Offset() to handle pagination.
func (rb *RecordBuilder) List(ctx context.Context) (*ResponseData, error) {
	if ctx == nil {
		ctx = rb.ctx
	}

	service := NewRecordService(rb.database, rb.layout, rb.client)
	// Note: Sorters are not yet supported in RecordBuilder.List, only in FindBuilder or direct Service calls.
	// You can add Sorter support to RecordBuilder if needed.
	return service.List(ctx, rb.offset, rb.limit)
}

// FindBuilder provides a fluent interface for building complex find operations.
type FindBuilder struct {
	client   *Client
	database string
	layout   string
	queries  []*groupQuery
	sorters  []*Sorter
	offset   string
	limit    string
	ctx      context.Context
}

// NewFindBuilder creates a new FindBuilder for the specified database and layout.
// Example:
//
//	builder := filemaker.NewFindBuilder(client, "Contacts", "ContactList")
//	response, err := builder.
//	    Where("Status", filemaker.Equal, "Active").
//	    Where("City", filemaker.Contains, "San").
//	    OrderBy("LastName", filemaker.Ascending).
//	    Limit(50).
//	    Execute(ctx)
func NewFindBuilder(client *Client, database, layout string) *FindBuilder {
	return &FindBuilder{
		client:   client,
		database: database,
		layout:   layout,
		queries:  make([]*groupQuery, 0),
		sorters:  make([]*Sorter, 0),
		ctx:      context.Background(),
	}
}

// WithContext sets the context for the operation.
func (fb *FindBuilder) WithContext(ctx context.Context) *FindBuilder {
	fb.ctx = ctx
	return fb
}

// Where adds a query condition to the current query group.
// Multiple Where calls are combined with AND logic.
// Use OrWhere to start a new query group (OR logic).
func (fb *FindBuilder) Where(field string, operator FieldOperator, value string) *FindBuilder {
	// If no query groups exist, create the first one
	if len(fb.queries) == 0 {
		fb.queries = append(fb.queries, NewGroupQuery())
	}

	// Add to the last query group
	lastGroup := fb.queries[len(fb.queries)-1]
	lastGroup.AddQuery(field, operator, value)

	return fb
}

// OrWhere starts a new query group with the specified condition.
// This creates an OR relationship with previous query groups.
func (fb *FindBuilder) OrWhere(field string, operator FieldOperator, value string) *FindBuilder {
	// Create a new query group
	newGroup := NewGroupQuery()
	newGroup.AddQuery(field, operator, value)
	fb.queries = append(fb.queries, newGroup)

	return fb
}

// OrderBy adds a sort order to the find operation.
func (fb *FindBuilder) OrderBy(field string, order SortOrder) *FindBuilder {
	fb.sorters = append(fb.sorters, NewSorter(field, order))
	return fb
}

// Offset sets the starting record for pagination.
func (fb *FindBuilder) Offset(offset int) *FindBuilder {
	fb.offset = strconv.Itoa(offset)
	return fb
}

// Limit sets the maximum number of records to return.
func (fb *FindBuilder) Limit(limit int) *FindBuilder {
	fb.limit = strconv.Itoa(limit)
	return fb
}

// Execute performs the find operation with the configured criteria.
func (fb *FindBuilder) Execute(ctx context.Context) (*ResponseData, error) {
	if ctx == nil {
		ctx = fb.ctx
	}

	service := NewSearchService(fb.database, fb.layout, fb.client)

	// Configure the search service
	service.GroupQueries(fb.queries...)

	if len(fb.sorters) > 0 {
		service.Sorters(fb.sorters...)
	}

	if fb.offset != "" {
		service.SetOffset(fb.offset)
	}

	if fb.limit != "" {
		service.SetLimit(fb.limit)
	}

	return service.Do(ctx)
}

// SessionBuilder provides a fluent interface for managing FileMaker sessions.
type SessionBuilder struct {
	client   *Client
	database string
	token    string
	username string
	password string
	ctx      context.Context
}

// NewSessionBuilder creates a new SessionBuilder for session management.
func NewSessionBuilder(client *Client, database string) *SessionBuilder {
	return &SessionBuilder{
		client:   client,
		database: database,
		ctx:      context.Background(),
	}
}

// WithContext sets the context for the operation.
func (sb *SessionBuilder) WithContext(ctx context.Context) *SessionBuilder {
	sb.ctx = ctx
	return sb
}

// WithCredentials sets the username and password for the session.
func (sb *SessionBuilder) WithCredentials(username, password string) *SessionBuilder {
	sb.username = username
	sb.password = password
	return sb
}

// Connect establishes a new session and returns the session token.
func (sb *SessionBuilder) Connect(ctx context.Context) (*ResponseData, error) {
	if sb.username != "" {
		return sb.CreateSession(ctx, WithBasicAuth(sb.username, sb.password))
	}
	return sb.CreateSession(ctx, sb.client.authProvider)
}

// ConnectWithDatasource establishes a session using external data source authentication.
func (sb *SessionBuilder) ConnectWithDatasource(ctx context.Context) (*ResponseData, error) {
	if sb.username != "" {
		return sb.CreateSession(ctx, WithCustomDatasource(sb.username, sb.password))
	}
	// Fallback to client default, assuming it is a datasource provider
	return sb.CreateSession(ctx, sb.client.authProvider)
}

// CreateSession establishes a new session using the provided authentication strategy.
func (sb *SessionBuilder) CreateSession(ctx context.Context, auth AuthProvider) (*ResponseData, error) {
	if ctx == nil {
		ctx = sb.ctx
	}

	response, err := sb.client.CreateSession(ctx, sb.database, auth)
	if err != nil {
		return nil, err
	}

	// Store the token for later operations
	if response != nil && response.Response.Token != "" {
		sb.token = response.Response.Token
	}

	return response, nil
}

// WithToken sets the session token for validation or disconnect operations.
func (sb *SessionBuilder) WithToken(token string) *SessionBuilder {
	sb.token = token
	return sb
}

// Validate checks if the current session token is still valid.
func (sb *SessionBuilder) Validate(ctx context.Context) (*ResponseData, error) {
	if ctx == nil {
		ctx = sb.ctx
	}

	if sb.token == "" {
		return nil, &ValidationError{
			Field:   "token",
			Message: "token is required for validation",
		}
	}

	return sb.client.ValidateSessionWithContext(ctx, sb.database, sb.token)
}

// Disconnect closes the session and invalidates the token.
func (sb *SessionBuilder) Disconnect(ctx context.Context) (*ResponseData, error) {
	if ctx == nil {
		ctx = sb.ctx
	}

	if sb.token == "" {
		return nil, &ValidationError{
			Field:   "token",
			Message: "token is required for disconnect",
		}
	}

	return sb.client.DisconnectWithContext(ctx, sb.database, sb.token)
}

// Token returns the current session token.
func (sb *SessionBuilder) Token() string {
	return sb.token
}
