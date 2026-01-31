# Global Fields

Global fields in FileMaker retain their values for the duration of a user's session. You can set these values using the `GlobalFieldsService`.

## Usage

Use `NewGlobalFieldsBuilder` to construct the map of fields and values, then send them using the service.

```go
// 1. Get a session token
session, _ := client.CreateSession(ctx, "MyDatabase")
token := session.Response.Token

// 2. Prepare the service
service := filemaker.NewGlobalFieldsService(client)

// 3. Build the fields map
fields := filemaker.NewGlobalFieldsBuilder().
    Add("MyGlobalTable::CurrentDate", "1/1/2024").
    Add("MyGlobalTable::CurrentUser", "Admin").
    Build()

// 4. Set the values
resp, err := service.SetGlobalFields(context.Background(), "MyDatabase", fields, token)

if err != nil {
    log.Fatal(err)
}
```
