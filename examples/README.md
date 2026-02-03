# FileMaker Go Library Examples

This directory contains executable examples demonstrating the core capabilities of the `filemaker` library. Each example is a self-contained Go program designed to run against a live FileMaker Server.

## Prerequisites

*   **FileMaker Server 18+** or FileMaker Cloud.
*   **A hosted database** with appropriate layouts and permissions (`fmrest` privilege).
*   **Go 1.25+** installed locally.

## Setup

The examples use environment variables to keep credentials secure and configurable. You must set these variables in your terminal before running the examples.

### Linux / macOS

```bash
export FM_HOST="https://your-filemaker-server.com"
export FM_DB="YourDatabaseName"
export FM_LAYOUT="YourLayoutName"
export FM_USER="admin"
export FM_PASS="password"

# Optional (for specific examples)
export FM_SCRIPT="TestScript"          # For scripts/main.go
export FM_EXT_USER="external-user"     # For authentication/main.go (ODBC/JDBC)
export FM_EXT_PASS="external-pass"
```

### Windows (PowerShell)

```powershell
$env:FM_HOST="https://your-filemaker-server.com"
$env:FM_DB="YourDatabaseName"
$env:FM_LAYOUT="YourLayoutName"
$env:FM_USER="admin"
$env:FM_PASS="password"
```

## Available Examples

### 1. Authentication
**Path:** `examples/authentication/main.go`
Demonstrates how to connect using:
*   **Basic Auth**: The standard connection method.
*   **External Data Sources**: Authenticating with separate credentials for ODBC/JDBC sources.

```bash
go run examples/authentication/main.go
```

### 2. Record Management (CRUD)
**Path:** `examples/records/main.go`
A complete walkthrough of the record lifecycle:
*   **Create** a new record.
*   **Read** (Get) the record by ID.
*   **Update** the record.
*   **Duplicate** the record.
*   **List** records with pagination (Limit/Offset).
*   **Delete** the record.

```bash
go run examples/records/main.go
```

### 3. Search (Find)
**Path:** `examples/search/main.go`
Shows how to use the `FindBuilder` for:
*   **Simple Finds**: Matching a single field.
*   **Complex Finds**: Using `OR` logic (multiple requests).
*   **Operators**: Using FileMaker operators (e.g., `>`, `==`).
*   **Sorting**: Ordering results.

```bash
go run examples/search/main.go
```

### 4. Scripts
**Path:** `examples/scripts/main.go`
Demonstrates script execution:
*   **Standalone**: Running a script directly via `ScriptService`.
*   **Triggers**: Attaching an `After-Request` script to a record creation event.

```bash
go run examples/scripts/main.go
```

### 5. Metadata
**Path:** `examples/metadata/main.go`
Shows how to inspect the server and solution schema:
*   Get **Product Information** (Server version).
*   List available **Databases**.
*   List **Layouts** in a database.
*   List **Scripts** in a database.

```bash
go run examples/metadata/main.go
```

### 6. Global Fields
**Path:** `examples/global_fields/main.go`
Demonstrates how to set session-scoped global field values.

```bash
go run examples/global_fields/main.go
```

### 7. Portal Pagination (v2.0+)
**Path:** `examples/portal_pagination/main.go`
Demonstrates advanced portal pagination:
*   **Portal Configuration**: Control offset and limit for each portal individually.
*   **SearchService**: Using `SetPortalConfigs()` for portal pagination.
*   **FindBuilder**: Using `WithPortals()` with fluent interface.
*   **Performance**: Reducing payload size with large portals.

```bash
go run examples/portal_pagination/main.go
```

**New Features in v2.0**:
- Portal-specific pagination (`_offset.{portal}`, `_limit.{portal}`)
- Fluent `PortalConfig` builder
- Better performance with large related datasets
