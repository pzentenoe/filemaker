# Search (Find)

The library provides a powerful `FindBuilder` and `SearchService` to perform complex queries against your FileMaker data.

## FindBuilder

The `FindBuilder` provides a fluent interface for constructing find requests, supporting multiple query groups (OR logic), sorting, offsets, limits, and portals.

### Basic Find

Find records where `FirstName` is "John".

```go
builder := filemaker.NewFindBuilder(client, "MyDatabase", "LayoutName")

resp, err := builder.
    AddQuery(filemaker.NewQueryField("FirstName", "John")).
    Do(context.Background())
```

### Multiple Criteria (AND)

Find records where `FirstName` is "John" AND `City` is "New York".

```go
q := filemaker.NewQuery().
    AddField("FirstName", "John").
    AddField("City", "New York")

resp, err := builder.AddQuery(q).Do(context.Background())
```

### Multiple Query Groups (OR)

Find records where (`FirstName` is "John") OR (`FirstName` is "Jane").

```go
resp, err := builder.
    AddQuery(filemaker.NewQueryField("FirstName", "John")).
    AddQuery(filemaker.NewQueryField("FirstName", "Jane")).
    Do(context.Background())
```

### Operators

You can use standard FileMaker find operators.

```go
q := filemaker.NewQuery().
    AddField("Price", ">100").
    AddField("CreatedDate", "1/1/2023...12/31/2023").
    AddField("Status", "==Active") // Exact match

resp, err := builder.AddQuery(q).Do(context.Background())
```

### Sorting, Limit, and Offset

```go
resp, err := builder.
    AddQuery(filemaker.NewQueryField("Status", "Active")).
    Sort("LastName", filemaker.Ascending).
    Sort("CreatedDate", filemaker.Descending).
    Limit(50).
    Offset(10).
    Do(context.Background())
```

### Portals

Specify which portals to return data for.

```go
resp, err := builder.
    AddQuery(filemaker.NewQueryField("ID", "123")).
    WithPortals("LineItems", "Notes").
    Do(context.Background())
```

## SearchService

If you prefer to build the `SearchData` structure manually or need lower-level control, you can use the `SearchService` directly, though `FindBuilder` is recommended for most use cases.

```go
service := filemaker.NewSearchService("MyDatabase", "LayoutName", client)

// ... configure service ...
resp, err := service.Do(context.Background())
```
