# Metadata

The Metadata Service allows you to discover information about the hosted databases, layouts, scripts, and server product info.

## Initialization

The Metadata Service is available via the `NewMetadataService` constructor or simply by using the client methods directly if you prefer wrapping them yourself, but the service struct is the standard way.

```go
metadataService := filemaker.NewMetadataService(client)
```

## Get Product Info

Retrieves information about the FileMaker Server, including version, build number, and date formats. This endpoint **does not** require a session token, but it does require valid Basic Auth credentials configured in the client.

```go
resp, err := metadataService.GetProductInfo(context.Background())
fmt.Printf("Server Version: %s\n", resp.Response.ProductInfo.Version)
```

## Get Databases

Retrieves a list of all databases hosted on the server that are accessible to the user defined in the client's Basic Auth. This endpoint **does not** require a session token.

```go
resp, err := metadataService.GetDatabases(context.Background())
for _, db := range resp.Response.Databases {
    fmt.Println(db.Name)
}
```

## Get Layouts

Retrieves a list of all layouts in a specific database. Requires a valid session token.

```go
// Get a token first
session, _ := client.CreateSession(ctx, "MyDatabase")
token := session.Response.Token

resp, err := metadataService.GetLayouts(context.Background(), "MyDatabase", token)
for _, layout := range resp.Response.Layouts {
    fmt.Println(layout.Name)
}
```

## Get Scripts

Retrieves a list of all scripts in a specific database. Requires a valid session token.

```go
resp, err := metadataService.GetScripts(context.Background(), "MyDatabase", token)
for _, script := range resp.Response.Scripts {
    fmt.Println(script.Name)
}
```

## Get Layout Metadata

Retrieves detailed information about a specific layout, including fields, value lists, and portals.

```go
resp, err := metadataService.GetLayoutMetadata(context.Background(), "MyDatabase", "LayoutName", token)
fmt.Printf("Field Count: %d\n", len(resp.Response.FieldMetaData))
```

