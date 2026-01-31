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
	user := os.Getenv("FM_USER")
	pass := os.Getenv("FM_PASS")

	if host == "" {
		return
	}

	client, _ := filemaker.NewClient(
		filemaker.SetURL(host),
		filemaker.SetBasicAuth(user, pass),
	)

	// 1. Simple Find
	// Find Active records
	fmt.Println("Searching for Active records...")
	builder := filemaker.NewFindBuilder(client, db, layout)

	resp, err := builder.
		Where("Status", filemaker.Equal, "Active").
		OrderBy("CreatedDate", filemaker.Descending).
		Limit(10).
		Execute(context.Background())

	if err != nil {
		// Handle "No records match" error specifically if needed
		log.Printf("Search failed: %v", err)
	} else {
		fmt.Printf("Found %d records\n", len(resp.Response.Data))
	}

	// 2. Complex Find (OR Search)
	// Status is Active OR Status is Pending
	fmt.Println("\nSearching for Active OR Pending records...")

	resp, err = filemaker.NewFindBuilder(client, db, layout).
		Where("Status", filemaker.Equal, "Active").
		OrWhere("Status", filemaker.Equal, "Pending").
		Execute(context.Background())

	if err == nil {
		fmt.Printf("Found %d records\n", len(resp.Response.Data))
	}

	// 3. Find with Operators
	// Price > 100 AND Category = Electronics
	fmt.Println("\nSearching for Expensive Electronics...")

	resp, err = filemaker.NewFindBuilder(client, db, layout).
		Where("Price", filemaker.GreaterThan, "100").
		Where("Category", filemaker.Equal, "Electronics").
		Execute(context.Background())

	if err == nil {
		fmt.Printf("Found %d records\n", len(resp.Response.Data))
	}
}
