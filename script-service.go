package filemaker

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	scriptExecutePath = "fmi/data/%s/databases/%s/layouts/%s/script/%s"
)

// ScriptService provides methods for executing FileMaker scripts.
type ScriptService interface {
	// Execute runs a FileMaker script with optional parameters.
	// The script parameter should be the name of the script to execute.
	// The scriptParam is passed to the script as a parameter (can be empty string).
	// Returns the script result in the response.
	Execute(ctx context.Context, database, layout, script, scriptParam, token string) (*ResponseData, error)

	// ExecuteAfterAction executes a script after a record operation.
	// This is typically used in conjunction with record operations (create, edit, delete).
	// Use the script.after query parameter for this functionality.
	ExecuteAfterAction(ctx context.Context, database, layout, recordID, script, scriptParam, token string, action string) (*ResponseData, error)
}

type scriptService struct {
	client *Client
}

// NewScriptService creates a new ScriptService instance.
func NewScriptService(client *Client) ScriptService {
	return &scriptService{
		client: client,
	}
}

// Execute runs a FileMaker script with optional parameters.
// The script is executed on the specified layout and requires a valid session token.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - database: Name of the FileMaker database
//   - layout: Name of the layout where the script will be executed
//   - script: Name of the script to execute
//   - scriptParam: Optional parameter to pass to the script (use empty string if none)
//   - token: Valid session token from Connect or ConnectWithDatasource
//
// The script result is returned in ResponseData.Response.ScriptResult field.
func (s *scriptService) Execute(ctx context.Context, database, layout, script, scriptParam, token string) (*ResponseData, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	if database == "" {
		return nil, &ValidationError{
			Field:   "database",
			Message: "database name is required",
		}
	}

	if layout == "" {
		return nil, &ValidationError{
			Field:   "layout",
			Message: "layout name is required",
		}
	}

	if script == "" {
		return nil, &ValidationError{
			Field:   "script",
			Message: "script name is required",
		}
	}

	if token == "" {
		return nil, &ValidationError{
			Field:   "token",
			Message: "session token is required",
		}
	}

	version := s.client.getVersion()

	// URL encode the script name to handle special characters
	encodedScript := url.PathEscape(script)
	path := fmt.Sprintf(scriptExecutePath, version, database, layout, encodedScript)

	headers := http.Header{}
	headers.Set("Authorization", "Bearer "+token)

	// Add script parameter as query parameter if provided
	params := url.Values{}
	if scriptParam != "" {
		params.Set("script.param", scriptParam)
	}

	options := &performRequestOptions{
		Method:  http.MethodGet,
		Path:    path,
		Params:  params,
		Headers: headers,
	}

	return s.client.executeQuery(ctx, options)
}

// ExecuteAfterAction executes a script after a record operation (create, edit, delete).
// This method is used when you want to run a script as part of a record modification workflow.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - database: Name of the FileMaker database
//   - layout: Name of the layout where the script will be executed
//   - recordID: ID of the record involved in the action (use "1" for create operations)
//   - script: Name of the script to execute after the action
//   - scriptParam: Optional parameter to pass to the script
//   - token: Valid session token
//   - action: The action to perform ("create", "edit", "delete")
//
// Note: This is a specialized method. For most use cases, use Execute() instead.
// For integrating scripts with record operations, use the script parameters in
// RecordService methods (e.g., Create with script.prerequest, script.presort).
func (s *scriptService) ExecuteAfterAction(ctx context.Context, database, layout, recordID, script, scriptParam, token string, action string) (*ResponseData, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	if database == "" {
		return nil, &ValidationError{
			Field:   "database",
			Message: "database name is required",
		}
	}

	if layout == "" {
		return nil, &ValidationError{
			Field:   "layout",
			Message: "layout name is required",
		}
	}

	if recordID == "" {
		return nil, &ValidationError{
			Field:   "recordID",
			Message: "record ID is required",
		}
	}

	if script == "" {
		return nil, &ValidationError{
			Field:   "script",
			Message: "script name is required",
		}
	}

	if token == "" {
		return nil, &ValidationError{
			Field:   "token",
			Message: "session token is required",
		}
	}

	if action == "" {
		return nil, &ValidationError{
			Field:   "action",
			Message: "action is required (create, edit, delete)",
		}
	}

	version := s.client.getVersion()

	var path string
	var method string
	var body any

	switch action {
	case "create":
		path = fmt.Sprintf("fmi/data/%s/databases/%s/layouts/%s/records", version, database, layout)
		method = http.MethodPost
		body = map[string]any{"fieldData": map[string]any{}}
	case "edit":
		path = fmt.Sprintf("fmi/data/%s/databases/%s/layouts/%s/records/%s", version, database, layout, recordID)
		method = http.MethodPatch
		body = map[string]any{"fieldData": map[string]any{}}
	case "delete":
		path = fmt.Sprintf("fmi/data/%s/databases/%s/layouts/%s/records/%s", version, database, layout, recordID)
		method = http.MethodDelete
		body = nil
	default:
		return nil, &ValidationError{
			Field:   "action",
			Message: "invalid action, must be create, edit, or delete",
		}
	}

	headers := http.Header{}
	headers.Set("Authorization", "Bearer "+token)

	// Add script execution parameters
	params := url.Values{}
	params.Set("script", script)
	if scriptParam != "" {
		params.Set("script.param", scriptParam)
	}

	options := &performRequestOptions{
		Method:      method,
		Path:        path,
		Params:      params,
		Headers:     headers,
		Body:        body,
		ContentType: "application/json",
	}

	return s.client.executeQuery(ctx, options)
}

// ScriptParameter represents script execution parameters that can be used
// with record operations (create, edit, delete, find).
type ScriptParameter struct {
	// Script name to execute
	Script string `json:"script,omitempty"`
	// Parameter to pass to the script
	ScriptParam string `json:"script.param,omitempty"`
}

// ToQueryParams converts ScriptParameter to URL query parameters.
func (sp *ScriptParameter) ToQueryParams(prefix string) url.Values {
	params := url.Values{}
	if sp.Script != "" {
		params.Set(prefix, sp.Script)
		if sp.ScriptParam != "" {
			params.Set(prefix+".param", sp.ScriptParam)
		}
	}
	return params
}

// MarshalJSON implements json.Marshaler for ScriptParameter.
func (sp *ScriptParameter) MarshalJSON() ([]byte, error) {
	m := make(map[string]string)
	if sp.Script != "" {
		m["script"] = sp.Script
		if sp.ScriptParam != "" {
			m["script.param"] = sp.ScriptParam
		}
	}
	return json.Marshal(m)
}

// ScriptContext holds script execution parameters for different stages
// of a FileMaker Data API request.
type ScriptContext struct {
	// PreRequest script executes before the action
	PreRequest *ScriptParameter `json:"-"`
	// PreSort script executes before sorting (for find operations)
	PreSort *ScriptParameter `json:"-"`
	// After script executes after the action completes
	After *ScriptParameter `json:"-"`
}

// ToQueryParams converts all script parameters to URL query parameters.
func (sc *ScriptContext) ToQueryParams() url.Values {
	params := url.Values{}

	if sc.PreRequest != nil && sc.PreRequest.Script != "" {
		for k, v := range sc.PreRequest.ToQueryParams("script.prerequest") {
			params[k] = v
		}
	}

	if sc.PreSort != nil && sc.PreSort.Script != "" {
		for k, v := range sc.PreSort.ToQueryParams("script.presort") {
			params[k] = v
		}
	}

	if sc.After != nil && sc.After.Script != "" {
		for k, v := range sc.After.ToQueryParams("script") {
			params[k] = v
		}
	}

	return params
}

// NewScriptContext creates a new ScriptContext with optional script parameters.
func NewScriptContext() *ScriptContext {
	return &ScriptContext{}
}

// WithPreRequest sets the pre-request script.
func (sc *ScriptContext) WithPreRequest(script, param string) *ScriptContext {
	sc.PreRequest = &ScriptParameter{
		Script:      script,
		ScriptParam: param,
	}
	return sc
}

// WithPreSort sets the pre-sort script (for find operations).
func (sc *ScriptContext) WithPreSort(script, param string) *ScriptContext {
	sc.PreSort = &ScriptParameter{
		Script:      script,
		ScriptParam: param,
	}
	return sc
}

// WithAfter sets the after script.
func (sc *ScriptContext) WithAfter(script, param string) *ScriptContext {
	sc.After = &ScriptParameter{
		Script:      script,
		ScriptParam: param,
	}
	return sc
}
