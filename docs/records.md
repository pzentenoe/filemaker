# Record Management

Performing CRUD (Create, Read, Update, Delete) operations is the core functionality of the FileMaker Data API. This library provides a fluent `RecordBuilder` to simplify these operations.

## RecordBuilder

The `RecordBuilder` is the recommended way to interact with records. It handles field data, portal data, scripts, and context configuration.

### Initialization

```go
builder := filemaker.NewRecordBuilder(client, "DatabaseName", "LayoutName")
```

## Create Record

```go
resp, err := filemaker.NewRecordBuilder(client, "Contacts", "ContactDetail").
    SetField("FirstName", "John").
    SetField("LastName", "Doe").
    SetField("Email", "john.doe@example.com").
    AddPortalRecord("PhoneNumbers", map[string]any{
        "Number": "555-0199",
        "Type":   "Mobile",
    }).
    Create(context.Background())

fmt.Println("New Record ID:", resp.Response.RecordID)
```

## Get Record

Retrieve a single record by its internal FileMaker Record ID.

```go
resp, err := filemaker.NewRecordBuilder(client, "Contacts", "ContactDetail").
    ForRecord("123").
    Get(context.Background())
```

## Update Record

Update specific fields of an existing record.

### Simple Update
```go
resp, err := filemaker.NewRecordBuilder(client, "Contacts", "ContactDetail").
    ForRecord("123").
    SetField("Status", "Active").
    Update(context.Background())
```

### Optimistic Locking (ModID)
To prevent overwriting changes made by others, provide the `ModId` (Modification ID) of the record version you are editing. The update will fail if the record has been modified since you retrieved it.

```go
resp, err := filemaker.NewRecordBuilder(client, "Contacts", "ContactDetail").
    ForRecord("123").
    WithModID("5"). // The ModId you received from a previous Get/Find
    SetField("Status", "Inactive").
    Update(context.Background())
```

## Delete Record

Delete a record by its ID.

```go
resp, err := filemaker.NewRecordBuilder(client, "Contacts", "ContactDetail").
    ForRecord("123").
    Delete(context.Background())
```

### Delete Related Records
You can also delete related records in a specific portal when deleting the parent record.

```go
resp, err := filemaker.NewRecordBuilder(client, "Contacts", "ContactDetail").
    ForRecord("123").
    WithDeleteRelated("PhoneNumbers"). // Portal name
    Delete(context.Background())
```

## Duplicate Record

Duplicate an existing record.

```go
resp, err := filemaker.NewRecordBuilder(client, "Contacts", "ContactDetail").
    ForRecord("123").
    Duplicate(context.Background())

fmt.Println("New Record ID:", resp.Response.RecordID)
```

## List Records

Retrieve a range of records.

```go
// Get 10 records starting from the first one
resp, err := filemaker.NewRecordBuilder(client, "Contacts", "ContactDetail").
    Limit(10).
    Offset(1).
    List(context.Background())
```