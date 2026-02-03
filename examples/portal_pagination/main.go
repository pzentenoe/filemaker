package main

import (
	"context"
	"fmt"
	"log"

	"github.com/pzentenoe/filemaker"
)

func main() {
	// Create FileMaker client
	client, err := filemaker.NewClient(
		filemaker.SetURL("https://your-server.com"),
		filemaker.SetBasicAuth("username", "password"),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Example 1: Using SearchService with portal pagination
	fmt.Println("Example 1: SearchService with portal pagination")
	searchWithPortalPagination(client)

	// Example 2: Using FindBuilder with portal pagination
	fmt.Println("\nExample 2: FindBuilder with portal pagination")
	findBuilderWithPortalPagination(client)
}

func searchWithPortalPagination(client *filemaker.Client) {
	// Create search service
	service := filemaker.NewSearchService("Contacts", "ContactList", client)

	// Create query
	query := filemaker.NewGroupQuery().
		AddQuery("Status", filemaker.Equal, "Active")

	// Configure portals with pagination
	portalConfigs := []*filemaker.PortalConfig{
		filemaker.NewPortalConfig("RelatedOrders").
			WithOffset(1).
			WithLimit(10),
		filemaker.NewPortalConfig("RelatedPayments").
			WithOffset(1).
			WithLimit(5),
	}

	// Execute search with portal pagination
	response, err := service.
		GroupQueries(query).
		SetPortalConfigs(portalConfigs...).
		Sorters(filemaker.NewSorter("LastName", filemaker.Ascending)).
		SetOffset("1").
		SetLimit("50").
		Do(context.Background())

	if err != nil {
		log.Printf("Search error: %v", err)
		return
	}

	fmt.Printf("Found %d records\n", len(response.Response.Data))

	// Process results with portal data
	for _, record := range response.Response.Data {
		fmt.Printf("Contact: %v\n", record.FieldData)

		// Access portal data (already []interface{} type)
		if orders, ok := record.PortalData["RelatedOrders"]; ok {
			fmt.Printf("  - Orders: %d records\n", len(orders))
		}

		if payments, ok := record.PortalData["RelatedPayments"]; ok {
			fmt.Printf("  - Payments: %d records\n", len(payments))
		}
	}
}

func findBuilderWithPortalPagination(client *filemaker.Client) {
	// Use FindBuilder with portal pagination
	response, err := filemaker.NewFindBuilder(client, "Contacts", "ContactList").
		Where("Status", filemaker.Equal, "Active").
		Where("City", filemaker.Contains, "San").
		WithPortals(
			filemaker.NewPortalConfig("RelatedOrders").
				WithOffset(1).
				WithLimit(10),
			filemaker.NewPortalConfig("RelatedPayments").
				WithOffset(1).
				WithLimit(5),
		).
		OrderBy("LastName", filemaker.Ascending).
		Offset(1).
		Limit(50).
		Execute(context.Background())

	if err != nil {
		log.Printf("Search error: %v", err)
		return
	}

	fmt.Printf("Found %d records\n", len(response.Response.Data))

	// Process results
	for _, record := range response.Response.Data {
		fmt.Printf("Contact: %v\n", record.FieldData)

		// Portal data is automatically included with pagination
		if orders, ok := record.PortalData["RelatedOrders"]; ok {
			fmt.Printf("  - Orders (paginated): %d records\n", len(orders))
		}

		if payments, ok := record.PortalData["RelatedPayments"]; ok {
			fmt.Printf("  - Payments (paginated): %d records\n", len(payments))
		}
	}
}
