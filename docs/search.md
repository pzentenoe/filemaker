# Search (Find)

The library provides a powerful `FindBuilder` and `SearchService` to perform complex queries against your FileMaker data.

## FindBuilder

The `FindBuilder` provides a fluent interface for constructing find requests, supporting multiple query groups (OR logic), sorting, offsets, limits, and portals.

### Basic Find

Find records where `FirstName` is "John".

```go
builder := filemaker.NewFindBuilder(client, "MyDatabase", "LayoutName")

resp, err := builder.
    Where("FirstName", filemaker.Equal, "John").
    Execute(context.Background())
```

### Multiple Criteria (AND)

Find records where `FirstName` is "John" AND `City` is "New York".

```go
resp, err := builder.
    Where("FirstName", filemaker.Equal, "John").
    Where("City", filemaker.Equal, "New York").
    Execute(context.Background())
```

### Multiple Query Groups (OR)

Find records where (`FirstName` is "John") OR (`FirstName` is "Jane").

```go
resp, err := builder.
    Where("FirstName", filemaker.Equal, "John").
    OrWhere("FirstName", filemaker.Equal, "Jane").
    Execute(context.Background())
```

### Operators

Use FileMaker field operators for advanced queries.

**Available Operators**:
- `Equal` - Exact match (`==`)
- `Contains` - Contains text (`==*value*`)
- `BeginsWith` - Starts with (`==value*`)
- `EndsWith` - Ends with (`==*value`)
- `GreaterThan` - Greater than (`>`)
- `GreaterThanEqual` - Greater than or equal (`>=`)
- `LessThan` - Less than (`<`)
- `LessThanEqual` - Less than or equal (`<=`)

```go
resp, err := builder.
    Where("Price", filemaker.GreaterThan, "100").
    Where("City", filemaker.Contains, "New").
    Where("Status", filemaker.Equal, "Active").
    Execute(context.Background())
```

### Sorting, Limit, and Offset

```go
resp, err := builder.
    Where("Status", filemaker.Equal, "Active").
    OrderBy("LastName", filemaker.Ascending).
    OrderBy("CreatedDate", filemaker.Descending).
    Limit(50).
    Offset(10).
    Execute(context.Background())
```

### Portals

#### Basic Portal Inclusion

Include portal data in the response:

```go
resp, err := builder.
    Where("ID", filemaker.Equal, "123").
    WithPortals(
        filemaker.NewPortalConfig("LineItems"),
        filemaker.NewPortalConfig("Notes"),
    ).
    Execute(context.Background())
```

#### Portal Pagination (v2.0+)

Control how many portal records are returned for each portal individually:

```go
resp, err := builder.
    Where("Status", filemaker.Equal, "Active").
    WithPortals(
        // Get first 10 line items
        filemaker.NewPortalConfig("LineItems").
            WithOffset(1).
            WithLimit(10),

        // Get first 5 notes
        filemaker.NewPortalConfig("Notes").
            WithOffset(1).
            WithLimit(5),
    ).
    Execute(context.Background())
```

**Benefits of Portal Pagination**:
- Reduces payload size for large portals
- Improves response time
- Enables lazy loading in UI
- Better resource utilization

**Query Parameters Generated**:
```
_offset.LineItems=1
_limit.LineItems=10
_offset.Notes=1
_limit.Notes=5
```

## SearchService

If you prefer to build the `SearchData` structure manually or need lower-level control, you can use the `SearchService` directly, though `FindBuilder` is recommended for most use cases.

```go
service := filemaker.NewSearchService("MyDatabase", "LayoutName", client)

// ... configure service ...
resp, err := service.Do(context.Background())
```
