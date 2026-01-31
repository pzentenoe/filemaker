package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pzentenoe/filemaker"
)

func main() {
	host := os.Getenv("FM_HOST")
	db := os.Getenv("FM_DB")
	layout := os.Getenv("FM_LAYOUT")
	scriptName := os.Getenv("FM_SCRIPT") // Name of a script to run
	user := os.Getenv("FM_USER")
	pass := os.Getenv("FM_PASS")

	if host == "" || scriptName == "" {
		return
	}

	client, _ := filemaker.NewClient(
		filemaker.SetURL(host),
		filemaker.SetBasicAuth(user, pass),
	)

	ctx := context.Background()

	// 1. Standalone Script Execution
	fmt.Println("Executing standalone script...")

	// Need a session token for standalone execution
	session, err := client.CreateSession(ctx, db)
	if err != nil {
		log.Fatal(err)
	}
	token := session.Response.Token
	defer client.Disconnect(db, token)

	scriptService := filemaker.NewScriptService(client)
	resp, err := scriptService.Execute(ctx, db, layout, scriptName, "some-param", token)
	if err != nil {
		log.Printf("Script execution failed: %v", err)
	} else {
		fmt.Printf("Script Result: %s\n", resp.Response.ScriptResult)
	}

	// 2. Script Triggered by Record Creation
	// This is often more efficient as it combines the action and script in one request
	fmt.Println("\nCreating record with script trigger...")

	createResp, err := filemaker.NewRecordBuilder(client, db, layout).
		SetField("FirstName", "Script").
		SetField("LastName", "Tester").
		// Run script AFTER the record is created
		WithAfterScript(scriptName, "triggered-by-create").
		Create(ctx)

	if err != nil {
		log.Printf("Create with script failed: %v", err)
	} else {
		fmt.Printf("Record created. Script Result: %s\n", createResp.Response.ScriptResult)

		// Cleanup
		_, _ = filemaker.NewRecordBuilder(client, db, layout).
			ForRecord(createResp.Response.RecordID).
			Delete(ctx)
	}
}
