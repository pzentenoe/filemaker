# Script Execution

You can execute FileMaker scripts either as standalone actions or attached to record operations (Create, Edit, Delete, Find).

## Script Service

For standalone script execution, use the `ScriptService`.

```go
scriptService := filemaker.NewScriptService(client)

// Execute a script
resp, err := scriptService.Execute(
    context.Background(),
    "MyDatabase",
    "LayoutName",
    "MyScriptName",
    "OptionalParameter",
    token, // Valid session token
)

if resp.Response.ScriptError != "0" {
    fmt.Printf("Script Error: %s\n", resp.Response.ScriptError)
}
fmt.Printf("Script Result: %s\n", resp.Response.ScriptResult)
```

## Scripts in Record Operations

The most common way to run scripts is by attaching them to record requests. This allows you to run scripts before the request is processed, before sorting (in finds), or after the request completes.

This is supported via the `RecordBuilder` and `FindBuilder`.

### RecordBuilder Example

```go
builder := filemaker.NewRecordBuilder(client, "MyDatabase", "LayoutName")

resp, err := builder.
    SetField("Status", "Pending").
    // Run "Validate" script before creating the record
    WithPreRequestScript("Validate", "param").
    // Run "LogAction" script after creating the record
    WithAfterScript("LogAction", "created").
    Create(context.Background())
```

### FindBuilder Example

```go
findBuilder := filemaker.NewFindBuilder(client, "MyDatabase", "LayoutName")

resp, err := findBuilder.
    AddQuery(filemaker.NewQueryField("Status", "Active")),
    // Run script before sorting the found set
    WithPreSortScript("CustomSortLogic", "").
    Do(context.Background())
```
