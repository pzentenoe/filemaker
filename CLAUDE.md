# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go library that provides a client interface for the FileMaker Data API. It enables Go applications to interact with FileMaker databases through RESTful API calls, handling authentication, CRUD operations, and complex search queries.

## Development Commands

### Testing
```bash
go test -v                    # Run all tests
go test -v -run TestName      # Run a specific test
```

### Build
```bash
go build                      # Build the library
go mod tidy                   # Clean up dependencies
```

## Architecture

### Core Components

**Client Layer** (client.go)
- `Client`: Central HTTP client that manages connections to FileMaker Server
- Handles all HTTP requests via `performRequest()` and `executeQuery()`
- Uses functional options pattern for configuration via `ClientOptions`
- Thread-safe with `sync.RWMutex` for concurrent access
- Supports custom HTTP clients and gzip compression

**Authentication** (auth.go)
- Session-based authentication via `Connect()` and `Disconnect()`
- Token-based authorization using Bearer tokens
- Supports external data source authentication with `ConnectWithDatasource()`
- Session validation via `ValidateSession()` to check if a token is still valid
- Sessions must be established before operations and disconnected after
- All authentication methods support context for cancellation and timeout

**Service Pattern**
The library uses a service-based architecture where each service handles a specific domain:

1. **RecordService** (query-service.go)
   - CRUD operations: Create, Edit, Delete, Duplicate, GetById, List
   - Each operation manages its own session lifecycle (Connect → Execute → Disconnect)
   - Uses `Payload` struct with `fieldData` and optional `portalData`
   - List operations support pagination (offset/limit) and sorting

2. **SearchService** (search-service.go)
   - Complex find operations using `/_find` endpoint
   - Builder pattern with method chaining: `GroupQueries()` → `Sorters()` → `Do()`
   - Supports multiple query groups (OR logic), sorting, pagination, and portal filtering
   - Query groups contain field operators (AND logic within group)

3. **MetadataService** (metadata-service.go)
   - Discovery and metadata operations for databases, layouts, and scripts
   - `GetDatabases()`: Lists all hosted databases accessible to authenticated user
   - `GetLayouts()`: Lists all layouts in a specified database
   - `GetLayoutMetadata()`: Retrieves detailed layout metadata (fields, value lists, portals, relationships)
   - `GetScripts()`: Lists all Data API-accessible scripts in a database
   - `GetProductInfo()`: Retrieves FileMaker Server product information (version, build, locale, date/time formats)
   - Methods requiring session token: GetLayouts, GetLayoutMetadata, GetScripts
   - Methods using basic auth only: GetDatabases, GetProductInfo

4. **ScriptService** (script-service.go)
   - Execute FileMaker scripts with parameters
   - `Execute()`: Runs a script on a specified layout with optional parameters
   - `ExecuteAfterAction()`: Runs a script after a record operation (create, edit, delete)
   - `ScriptContext`: Builder pattern for managing pre-request, pre-sort, and after scripts
   - `ScriptParameter`: Encapsulates script name and parameter for reusable script execution
   - Scripts can be integrated with record operations using query parameters
   - Script results returned in ResponseData.Response.ScriptResult field

5. **GlobalFieldsService** (global-fields-service.go)
   - Set global field values for the current database session
   - `SetGlobalFields()`: Sets one or more global fields using a map of field names to values
   - `GlobalFieldsBuilder`: Fluent builder pattern for constructing global field maps
   - `GlobalField`: Helper type for programmatic global field creation
   - Global fields maintain their values across all records during a session
   - Useful for user preferences, context, and session-specific data

6. **ContainerService** (container-service.go)
   - Upload files to container fields in FileMaker databases
   - `UploadFile()`: Uploads a file from the local filesystem
   - `UploadFileWithRepetition()`: Uploads to a specific repetition of a repeating container field
   - `UploadData()`: Uploads binary data from memory
   - `UploadDataWithRepetition()`: Uploads binary data to a specific repetition
   - Supports multipart form uploads for images, PDFs, and other file types
   - `ContainerFileInfo`: Helper type for organizing file upload information
   - Handles repeating fields with 1-based indexing

**Query Building Components**

- **queryFieldOperator** (query-field.go): Represents a single field query with operator
  - Operators: Equal, Contains, BeginsWith, EndsWith, GreaterThan, GreaterThanEqual, LessThan, LessThanEqual
  - Converts to FileMaker query syntax (e.g., `==value`, `==*value*`, `>value`)

- **groupQuery** (query-group.go): Groups multiple field operators (AND logic)
  - Multiple groups in SearchService create OR conditions

- **Sorter** (sorter.go): Defines sort order for results
  - Supports Ascending/Descending order
  - Can be used in both List and Search operations

### API Path Structure

All API paths follow the pattern:
```
fmi/data/{version}/databases/{database}/layouts/{layout}/{operation}
```

Where `{version}` defaults to "vLatest" (configurable via `SetVersion()`)

### Request/Response Flow

1. Service method called → `Connect()` to get session token
2. Build `performRequestOptions` with method, path, body, headers
3. Execute via `client.executeQuery()`
4. Parse response into `ResponseData` structure
5. `Disconnect()` via defer to clean up session
6. Return `ResponseData` with response data and messages

### Data Structures

**ResponseData** (search-data.go)
- Top-level response wrapper containing `Response` and `Messages`
- `Response` includes token, recordID, modID, dataInfo, and data array
- `Datum` represents individual records with fieldData, portalData, recordID, modID
- `DataInfo` provides metadata: database, layout, table, record counts

## Key Patterns

1. **Functional Options Pattern**: Client configuration uses option functions (`SetURL()`, `SetUsername()`, etc.)

2. **Service Factories**: Services created via `NewRecordService()` and `NewSearchService()` with database, layout, and client

3. **Session Management**: Each service operation handles its own session lifecycle with defer for cleanup

4. **Builder Pattern**: SearchService uses method chaining for complex queries

5. **Thread Safety**: Client uses RWMutex for safe concurrent access during authentication operations

## Important Implementation Notes

- All service operations automatically manage FileMaker sessions (Connect/Disconnect)
- The client supports custom HTTP clients for advanced configurations (e.g., TLS, proxy, timeouts)
- Portal data is optional in all operations via the `omitempty` tag
- Query operators are converted to FileMaker-specific syntax in `valueWithOp()`
- Sorting can be applied to both List and Search operations
- The library requires Go 1.23+ (uses modern features like generics, `any` type, and improved error handling)

## Recent Improvements (v2.0)

### Context Support
All service methods now support `context.Context` for cancellation and timeout control:
- **Auth**: `ConnectWithContext()`, `DisconnectWithContext()`, `ConnectWithDatasourceContext()`, `ValidateSessionWithContext()`
- **RecordService**: All methods (Create, Edit, Delete, GetById, List, Duplicate) accept `ctx context.Context` as first parameter
- **SearchService**: `Do(ctx context.Context)` method accepts context
- **Session Validation**: `ValidateSession()` checks if a token is still valid without performing operations
- Original methods without context remain for backward compatibility and use `context.Background()` internally

### Error Handling
Comprehensive error system with typed errors:
- **FileMakerError**: Structured errors with code, message, and HTTP status
- **ValidationError**: Input validation failures
- **AuthenticationError**: Authentication failures
- **NetworkError**: Network-level errors
- **TimeoutError**: Timeout scenarios
- Error helpers: `IsRetryable()`, `IsAuthError()`, `ParseFileMakerError()`
- Proper error wrapping with `%w` verb for error chains

### Retry Logic and Resilience
Automatic retry with exponential backoff for transient failures:
- **RetryConfig**: Configurable retry behavior (max retries, wait times, jitter)
- Default retries: 3 attempts with exponential backoff (1s to 30s)
- Retryable status codes: 429 (rate limit), 500, 502, 503, 504
- Smart retry detection using `IsRetryable()` helper for FileMaker error codes (952, 953)
- Context-aware retries with cancellation support
- Optional `OnRetry` callback for logging/monitoring
- Configuration options: `SetRetryConfig()`, `SetMaxRetries()`, `SetRetryWaitTime()`, `DisableRetry()`

### HTTP Client Configuration
Comprehensive timeout and connection management:
- **HTTPClientConfig**: Fine-grained control over HTTP behavior
- Configurable timeouts: overall request, dial, TLS handshake, response headers, idle connections
- Connection pooling: max idle connections, max connections per host
- Optional TLS configuration and compression settings
- Configuration options: `SetTimeout()`, `SetHTTPClientConfig()`, `SetDialTimeout()`, `SetTLSHandshakeTimeout()`
- Defaults optimized for FileMaker Server (30s timeout, 10s dial/TLS, HTTP/2 enabled)

### Metadata Discovery (NEW)
Complete metadata and discovery capabilities:
- **MetadataService**: New service for discovering databases, layouts, and scripts
- Database discovery: List all hosted databases accessible to user
- Layout operations: List layouts and retrieve detailed metadata (fields, value lists, portals)
- Script discovery: List all Data API-accessible scripts
- Product information: Get FileMaker Server version, build, and configuration
- Full context support for all operations

### Script Execution (NEW)
Comprehensive script execution capabilities:
- **ScriptService**: Execute FileMaker scripts with parameters
- Direct script execution with `Execute()` method
- Script execution after record operations with `ExecuteAfterAction()`
- `ScriptContext`: Builder pattern for managing pre-request, pre-sort, and after scripts in workflows
- `ScriptParameter`: Reusable script parameter encapsulation
- Script results available in ResponseData.Response.ScriptResult
- URL-safe script name encoding for special characters
- Full context support and validation

### Global Fields Management (NEW)
Session-wide global field management:
- **GlobalFieldsService**: Set global field values for the current database session
- Global fields maintain values across all records during a session
- `GlobalFieldsBuilder`: Fluent builder pattern for constructing global field maps
- Useful for user preferences, session context, and application state
- Full validation and context support

### Container Fields (NEW)
File upload capabilities for container fields:
- **ContainerService**: Upload files to FileMaker container fields
- Upload from filesystem with `UploadFile()` or from memory with `UploadData()`
- Support for repeating container fields with 1-based indexing
- Multipart form data uploads for images, PDFs, and other file types
- `ContainerFileInfo`: Helper type for organizing file upload metadata
- Retry logic and error handling for reliable uploads
- Full context support for cancellation and timeouts

### Fluent Builder Pattern (NEW)
Improved developer experience with fluent interfaces:
- **RecordBuilder**: Fluent API for CRUD operations with method chaining
  - `SetField()`, `SetFields()`: Configure field data
  - `AddPortalRecord()`, `SetPortalRecords()`: Manage portal data
  - `WithScripts()`, `WithPreRequestScript()`, `WithAfterScript()`: Script integration
  - `ForRecord()`: Set record ID for updates/deletes
  - `Create()`, `Update()`, `Delete()`, `Get()`, `Duplicate()`: Execute operations
- **FindBuilder**: Simplified search with SQL-like syntax
  - `Where()`: Add AND conditions
  - `OrWhere()`: Add OR conditions (new query group)
  - `OrderBy()`: Sort results
  - `Offset()`, `Limit()`: Pagination
  - `Execute()`: Perform the search
- **SessionBuilder**: Session management with fluent interface
  - `Connect()`, `ConnectWithDatasource()`: Establish sessions
  - `WithToken()`: Set session token
  - `Validate()`: Check token validity
  - `Disconnect()`: Close session
  - `Token()`: Get current token
- All builders support `WithContext()` for context management

### Logging and Observability (NEW)
Optional logging and metrics for production monitoring:
- **Logger Interface**: Pluggable logging with Debug, Info, Warn, Error levels
  - `NoOpLogger`: Default (no logging)
  - `StandardLogger`: Built-in logger writing to stdout
  - Custom loggers via Logger interface
- **Configuration**: `SetLogger()`, `EnableLogging(level)`
- **Metrics**: Operational metrics tracking
  - Request counts (total, succeeded, failed)
  - Retry attempts
  - Session lifecycle (created, closed, active)
  - `Metrics.Snapshot()`: Get current metrics state
- **Request Tracing**: RequestContext for correlation
  - Request IDs for distributed tracing
  - Operation metadata (database, layout, duration)
  - Context propagation with `WithRequestContext()`
- **Configuration**: `SetMetrics()`, `EnableMetrics()`

### Modern Go Features
- Uses `any` instead of `interface{}`
- Uses `io.ReadAll` and `io.NopCloser` instead of deprecated `ioutil` package
- Error wrapping with fmt.Errorf and %w verb
- Thread-safe operations with proper lock management (read locks released before I/O)
- Atomic operations for metrics (sync/atomic)