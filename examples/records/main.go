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

	if host == "" || layout == "" {
		fmt.Println("Please set FM_HOST, FM_DB, FM_LAYOUT, FM_USER, FM_PASS")
		return
	}

	client, _ := filemaker.NewClient(
		filemaker.SetURL(host),
		filemaker.SetBasicAuth(user, pass),
	)

	ctx := context.Background()
	builder := filemaker.NewRecordBuilder(client, db, layout)

	// 1. Create a Record
	fmt.Println("Creating record...")
	createResp, err := builder.
		SetField("FirstName", "John").
		SetField("LastName", "Doe").
		SetField("Status", "Active").
		Create(ctx)

	if err != nil {
		log.Fatalf("Create failed: %v", err)
	}
	recordID := createResp.Response.RecordID
	fmt.Printf("Created Record ID: %s\n", recordID)

	// 2. Get the Record
	fmt.Println("\nReading record...")
	getResp, err := filemaker.NewRecordBuilder(client, db, layout).
		ForRecord(recordID).
		Get(ctx)

	if err != nil {
		log.Printf("Get failed: %v", err)
	} else {
		// Access field data
		if len(getResp.Response.Data) > 0 {
			fields := getResp.Response.Data[0].FieldData
			fmt.Printf("Name: %v %v\n", fields["FirstName"], fields["LastName"])
		}
	}

	// 3. Update the Record
	fmt.Println("\nUpdating record...")
	updateResp, err := filemaker.NewRecordBuilder(client, db, layout).
		ForRecord(recordID).
		SetField("Status", "Inactive").
		Update(ctx)

	if err != nil {
		log.Printf("Update failed: %v", err)
	} else {
		fmt.Printf("Record updated. New ModID: %s\n", updateResp.Response.ModID)
	}

	// 4. Duplicate the Record
	fmt.Println("\nDuplicating record...")
	dulpResp, err := filemaker.NewRecordBuilder(client, db, layout).
		ForRecord(recordID).
		Duplicate(ctx)

	if err != nil {
		log.Printf("Duplicate failed: %v", err)
	} else {
		fmt.Printf("Duplicated Record ID: %s\n", dulpResp.Response.RecordID)
		// Clean up duplicate
		_, _ = filemaker.NewRecordBuilder(client, db, layout).
			ForRecord(dulpResp.Response.RecordID).
			Delete(ctx)
	}

	// 5. Delete the Original Record
	fmt.Println("\nDeleting original record...")
	_, err = filemaker.NewRecordBuilder(client, db, layout).
		ForRecord(recordID).
		Delete(ctx)

	if err != nil {
		log.Printf("Delete failed: %v", err)
	} else {
		fmt.Println("Record deleted successfully.")
	}

	// 6. List Records
	fmt.Println("\nListing records...")
	listResp, err := filemaker.NewRecordBuilder(client, db, layout).
		Limit(5).
		Offset(1).
		List(ctx)

	if err != nil {
		log.Printf("List failed: %v", err)
	} else {
		fmt.Printf("Retrieved %d records\n", len(listResp.Response.Data))
	}
}
