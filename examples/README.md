# FileMaker Go Library Examples

This directory contains executable examples demonstrating how to use the `filemaker` library.

## Prerequisites

To run these examples, you need access to a FileMaker Server (v18+) and a database file hosted on it.

## Environment Variables

The examples rely on environment variables to configure the connection. Set these before running the examples:

```bash
export FM_HOST="https://your-server.com"
export FM_DB="YourDatabaseName"
export FM_LAYOUT="YourLayoutName"
export FM_USER="admin"
export FM_PASS="password"
# Optional
export FM_SCRIPT="YourScriptName" 
export FM_EXT_USER="external-user" # For datasource auth example
export FM_EXT_PASS="external-pass"
```

## Running Examples

Navigate to the specific example folder and run `go run main.go`.

### Authentication
Demonstrates Basic Auth and External Data Source authentication.
```bash
go run examples/authentication/main.go
```

### Records
Demonstrates CRUD operations (Create, Get, Update, Delete) and Listing.
```bash
go run examples/records/main.go
```

### Search
Demonstrates simple and complex (AND/OR) find requests.
```bash
go run examples/search/main.go
```

### Scripts
Demonstrates executing scripts standalone and as triggers during record operations.
```bash
go run examples/scripts/main.go
```

### Metadata
Demonstrates fetching server info, database lists, layouts, and scripts.
```bash
go run examples/metadata/main.go
```

### Global Fields
Demonstrates setting session-scoped global field values.
```bash
go run examples/global_fields/main.go
```
