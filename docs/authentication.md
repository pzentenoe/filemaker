# Authentication

The FileMaker Data API library supports multiple authentication strategies, allowing you to connect to FileMaker Server or FileMaker Cloud using the most appropriate method for your security requirements.

## Overview

Authentication is managed via the `CreateSession` method on the `Client`. The library provides several built-in strategies:

*   **Basic Auth**: Standard username and password for the FileMaker file.
*   **OAuth**: Authenticate using an external Identity Provider (IdP).
*   **Claris ID**: Authenticate using a Claris ID token (FMID) for FileMaker Cloud.
*   **Data Source**: Authenticate against an external ODBC/JDBC data source while authenticating to the master file.

## Basic Authentication

This is the most common method for on-premise FileMaker Server.

```go
client, _ := filemaker.NewClient(
    filemaker.SetURL("https://your-server.com"),
    filemaker.SetBasicAuth("admin", "password"),
)

// Create a session using the default Basic Auth provider configured in NewClient
resp, err := client.CreateSession(context.Background(), "MyDatabase")

// Or explicitly with a specific provider
resp, err := client.CreateSession(context.Background(), "MyDatabase", filemaker.WithBasicAuth("admin", "password"))
```

## OAuth Authentication

For solutions configured with external identity providers (like Microsoft Azure AD, Google, etc.).

```go
client, _ := filemaker.NewClient(
    filemaker.SetURL("https://your-server.com"),
)

// You need the Request ID and Identifier from your OAuth flow
resp, err := client.CreateSession(
    context.Background(),
    "MyDatabase",
    filemaker.WithOAuth("request-id-123", "identifier-456"),
)
```

## Claris ID (FMID) Authentication

For FileMaker Cloud solutions.

```go
client, _ := filemaker.NewClient(
    filemaker.SetURL("https://your-server.com"),
)

resp, err := client.CreateSession(
    context.Background(),
    "MyDatabase",
    filemaker.WithFMID("your-claris-id-token"),
)
```

## External Data Source Authentication

To authenticate with credentials for an external ODBC/JDBC data source configured within your FileMaker file.

**Note:** This strategy uses the client's configured Basic Auth credentials to authenticate with the **Master File**, and the provided credentials for the **External Data Source**.

```go
// 1. Configure client with Master File credentials
client, _ := filemaker.NewClient(
    filemaker.SetURL("https://your-server.com"),
    filemaker.SetBasicAuth("master-file-user", "master-file-pass"),
)

// 2. Connect providing External Data Source credentials
resp, err := client.CreateSession(
    context.Background(),
    "MyDatabase",
    filemaker.WithCustomDatasource("external-db-user", "external-db-pass"),
)
```

## Session Management

### Validate Session

Check if a token is still valid:

```go
resp, err := client.ValidateSession("MyDatabase", "your-session-token")
```

### Disconnect (Logout)

Always disconnect when done to free up user licenses on the server.

```go
_, err := client.Disconnect("MyDatabase", "your-session-token")
```
