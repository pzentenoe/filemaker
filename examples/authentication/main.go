package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pzentenoe/filemaker"
)

func main() {
	// Configuration from environment variables
	host := os.Getenv("FM_HOST")
	db := os.Getenv("FM_DB")
	user := os.Getenv("FM_USER")
	pass := os.Getenv("FM_PASS")

	if host == "" {
		fmt.Println("Please set FM_HOST, FM_DB, FM_USER, FM_PASS environment variables")
		return
	}

	// 1. Basic Authentication (Most Common)
	client, err := filemaker.NewClient(
		filemaker.SetURL(host),
		filemaker.SetBasicAuth(user, pass),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Create a session explicitly
	session, err := client.CreateSession(context.Background(), db)
	if err != nil {
		log.Printf("Failed to create session: %v", err)
	} else {
		fmt.Printf("Session created successfully. Token: %s\n", session.Response.Token)

		// Logout
		_, _ = client.Disconnect(db, session.Response.Token)
	}

	// 2. External Data Source Authentication
	// Used when you need to authenticate to an external ODBC/JDBC source
	// defined in your FileMaker solution.
	extUser := os.Getenv("FM_EXT_USER")
	extPass := os.Getenv("FM_EXT_PASS")

	if extUser != "" {
		fmt.Println("\nAttempting External Data Source Auth...")

		// Note: The client still needs the master file credentials (set above via NewClient)
		dsSession, err := client.CreateSession(
			context.Background(),
			db,
			filemaker.WithCustomDatasource(extUser, extPass),
		)

		if err != nil {
			log.Printf("Failed to create datasource session: %v", err)
		} else {
			fmt.Printf("Datasource session created. Token: %s\n", dsSession.Response.Token)
			_, _ = client.Disconnect(db, dsSession.Response.Token)
		}
	}
}
