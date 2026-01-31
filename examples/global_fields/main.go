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
	user := os.Getenv("FM_USER")
	pass := os.Getenv("FM_PASS")

	if host == "" {
		return
	}

	client, _ := filemaker.NewClient(
		filemaker.SetURL(host),
		filemaker.SetBasicAuth(user, pass),
	)

	// 1. Create Session
	session, err := client.CreateSession(context.Background(), db)
	if err != nil {
		log.Fatal(err)
	}
	token := session.Response.Token
	defer func(client *filemaker.Client, database, token string) {
		_, _ = client.Disconnect(database, token)
	}(client, db, token)

	// 2. Build Global Fields Map
	// Ensure these fields are actually global fields in your solution
	// Format: TableName::FieldName
	fields := filemaker.NewGlobalFieldsBuilder().
		Add("Globals::CurrentLanguage", "en").
		Add("Globals::UserRole", "Admin").
		Build()

	// 3. Set Values
	fmt.Println("Setting global fields...")
	service := filemaker.NewGlobalFieldsService(client)
	_, err = service.SetGlobalFields(context.Background(), db, fields, token)

	if err != nil {
		log.Printf("Failed to set globals: %v", err)
	} else {
		fmt.Println("Global fields set successfully.")
	}
}
