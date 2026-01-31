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
	user := os.Getenv("FM_USER")
	pass := os.Getenv("FM_PASS")

	if host == "" {
		return
	}

	client, _ := filemaker.NewClient(
		filemaker.SetURL(host),
		filemaker.SetBasicAuth(user, pass),
	)

	service := filemaker.NewMetadataService(client)
	ctx := context.Background()

	// 1. Get Product Info (No session required)
	fmt.Println("Fetching Product Info...")
	info, err := service.GetProductInfo(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Server: %s %s\n", info.Response.ProductInfo.Name, info.Response.ProductInfo.Version)

	// 2. Get Databases (No session required)
	fmt.Println("\nFetching Databases...")
	dbs, err := service.GetDatabases(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for _, db := range dbs.Response.Databases {
		fmt.Printf("- %s\n", db.Name)
	}

	// 3. Get Layouts (Session required)
	// We'll pick the first database found or use env var
	dbName := os.Getenv("FM_DB")
	if dbName == "" && len(dbs.Response.Databases) > 0 {
		dbName = dbs.Response.Databases[0].Name
	}

	if dbName != "" {
		fmt.Printf("\nFetching Layouts for %s...\n", dbName)

		// Create session
		session, err := client.CreateSession(ctx, dbName)
		if err != nil {
			log.Fatal(err)
		}
		token := session.Response.Token
		defer client.Disconnect(dbName, token)

		layouts, err := service.GetLayouts(ctx, dbName, token)
		if err != nil {
			log.Printf("Failed to get layouts: %v", err)
		} else {
			for _, l := range layouts.Response.Layouts {
				fmt.Printf("- %s\n", l.Name)
			}
		}

		// 4. Get Scripts
		fmt.Printf("\nFetching Scripts for %s...\n", dbName)
		scripts, err := service.GetScripts(ctx, dbName, token)
		if err != nil {
			log.Printf("Failed to get scripts: %v", err)
		} else {
			for _, s := range scripts.Response.Scripts {
				fmt.Printf("- %s\n", s.Name)
			}
		}
	}
}
