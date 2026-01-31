# FileMaker Go Library

![CI](https://github.com/pzentenoe/filemaker/actions/workflows/actions.yml/badge.svg)
[![codecov](https://codecov.io/github/pzentenoe/filemaker/graph/badge.svg?token=3W164MZ18S)](https://codecov.io/github/pzentenoe/filemaker)
[![Go Report Card](https://goreportcard.com/badge/github.com/pzentenoe/filemaker)](https://goreportcard.com/report/github.com/pzentenoe/filemaker)
![License](https://img.shields.io/github/license/pzentenoe/filemaker)
![GitHub release](https://img.shields.io/github/v/release/pzentenoe/filemaker)

[![FOSSA Status](https://app.fossa.com/api/projects/custom%2B60042%2Fgit%40github.com%3Apzentenoe%2Ffilemaker.git.svg?type=shield&issueType=license)](https://app.fossa.com/projects/custom%2B60042%2Fgit%40github.com%3Apzentenoe%2Ffilemaker.git?ref=badge_shield&issueType=license)
[![FOSSA Status](https://app.fossa.com/api/projects/custom%2B60042%2Fgit%40github.com%3Apzentenoe%2Ffilemaker.git.svg?type=shield&issueType=security)](https://app.fossa.com/projects/custom%2B60042%2Fgit%40github.com%3Apzentenoe%2Ffilemaker.git?ref=badge_shield&issueType=security)

**A robust, idiomatic, and type-safe Go client for the FileMaker Data API (v18+).**

Welcome! This library allows you to interact with your FileMaker Server or Cloud databases using the Go programming language. It is designed to be developer-friendly, abstracting away the raw HTTP requests into a clean, fluent API.

[Official FileMaker Data API Documentation](https://fmhelp.filemaker.com/docs/18/es/dataapi/)

## Key Features

*   **🔐 Authentication**: Seamless support for:
    *   **Basic Auth** (Standard FileMaker Security)
    *   **OAuth** (External Identity Providers)
    *   **Claris ID** (FileMaker Cloud)
    *   **External Data Sources** (ODBC/JDBC integration)
*   **💾 Record Management**: Fluent `RecordBuilder` for easy CRUD operations:
    *   **Create**, **Read**, **Update**, **Delete**
    *   **Duplicate** records
    *   **List** records with pagination
    *   **Portal** data manipulation
    *   **Optimistic Locking** via `modId`
*   **🔍 Advanced Search**: Powerful `FindBuilder` supporting:
    *   Complex queries with **AND** / **OR** logic
    *   FileMaker find operators (e.g., `==`, `...`, `>`, etc.)
    *   Sorting and Pagination (`Limit`, `Offset`)
*   **📜 Script Execution**:
    *   Run standalone scripts
    *   Trigger scripts **Pre-Request**, **Pre-Sort**, or **After-Request** during record operations
*   **📁 Metadata Discovery**:
    *   Inspect **Databases**, **Layouts**, and **Scripts**
    *   Retrieve **Product Information**
*   **🌐 Global Fields**: Set session-scoped global field values easily.
*   **⚙️ Robustness**:
    *   Context-aware (`context.Context`) for timeouts and cancellation
    *   Configurable **Retry Logic** for network resilience
    *   Comprehensive Error Handling

## Prerequisites

Before you begin, ensure you have the following:

*   **Go**: Version **1.25** or higher.
*   **FileMaker Server**: Version 18 or newer (or FileMaker Cloud).
*   **Access**: A user account on your FileMaker database with the `fmrest` extended privilege enabled.

## Installation

To install the library, run the following command in your terminal:

```bash
go get github.com/pzentenoe/filemaker
```

Then, import it in your Go code:

```go
import "github.com/pzentenoe/filemaker"
```

## Usage

Here is a quick example to get you up and running. This demonstrates connecting to a server and creating a new record.

```go
package main

import (
	"context"
	"fmt"
	"log"
	
	"github.com/pzentenoe/filemaker"
)

func main() {
	// 1. Initialize the Client
	// Configure the client with your server details and credentials.
	client, err := filemaker.NewClient(
		filemaker.SetURL("https://your-fm-server.com"),
		filemaker.SetBasicAuth("admin", "password"),
		// Optional: Configure timeouts or retry logic
		// filemaker.SetTimeout(30 * time.Second),
	)
	if err != nil {
		log.Fatalf("Failed to initialize client: %v", err)
	}

	// 2. Create a Record using the Fluent Builder
	// The builder interface makes it easy to construct requests.
	// Note: Session management (login/logout) is handled automatically for this request.
	resp, err := filemaker.NewRecordBuilder(client, "ContactsDB", "ContactLayout").
		SetField("FirstName", "John").
		SetField("LastName", "Doe").
		SetField("Email", "john.doe@example.com").
		Create(context.Background())

	if err != nil {
		log.Fatalf("Error creating record: %v", err)
	}

	fmt.Printf("Successfully created record with ID: %s\n", resp.Response.RecordID)
}
```

## Documentation

We have comprehensive documentation available in the `docs/` directory to help you master every feature:

| Topic | Description |
| :--- | :--- |
| **[Authentication](docs/authentication.md)** | Connecting with Basic Auth, OAuth, Claris ID, and External Data Sources |
| **[Record Management](docs/records.md)** | Performing CRUD operations, managing Portals, and Listing records |
| **[Search (Find)](docs/search.md)** | Building complex find requests, using operators, and sorting results |
| **[Scripts](docs/scripts.md)** | Executing FileMaker scripts and handling triggers |
| **[Metadata](docs/metadata.md)** | Discovering database schemas, layouts, and scripts |
| **[Global Fields](docs/global_fields.md)** | Managing global field values |

## Examples

Check out the `examples/` directory for complete, runnable code snippets covering all major features:

*   [Authentication Examples](examples/authentication/main.go)
*   [CRUD & List Examples](examples/records/main.go)
*   [Search Examples](examples/search/main.go)
*   [Script Examples](examples/scripts/main.go)
*   [Metadata Examples](examples/metadata/main.go)

## Contributing

Contributions are welcome! If you find a bug or want to request a feature, please open an issue or submit a pull request.

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for a history of changes to this project.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

Copyright (c) 2026 Pablo Zenteno

## Author

*   **[Pablo Zenteno](https://github.com/pzentenoe)**

---
<a href="https://www.buymeacoffee.com/pzentenoe" target="_blank"><img src="https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png" alt="Buy Me A Coffee" style="height: 41px !important;width: 174px !important;box-shadow: 0px 3px 2px 0px rgba(190, 190, 190, 0.5) !important;-webkit-box-shadow: 0px 3px 2px 0px rgba(190, 190, 190, 0.5) !important;" ></a>