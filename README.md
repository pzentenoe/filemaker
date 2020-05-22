# Filemaker

**This is a library created in go to connect to filemaker.**

[Filemaker API](https://fmhelp.filemaker.com/docs/18/es/dataapi/) documentation.

**Example:**


To use the required version of Elastic in your application, you
should use [Go modules](https://github.com/golang/go/wiki/Modules)
to manage dependencies. Make sure to use a version such as `1.0.0` or later.

To use Elastic, import:

```go
import "github.com/pzentenoe/filemaker"
```

## Getting Started

The first thing you do is to create a [Client](https://github.com/pzentenoe/filemaker/blob/master/client.go).
The client connects to Filemaker passing a host.


```go
package main

import (
	"fmt"
	"github.com/pzentenoe/filemaker"
	"net/http"
)

func main() {

	client, err := filemaker.NewClient(
		filemaker.SetURL("https://localhost:port"),
		filemaker.SetUsername("username"),
		filemaker.SetPassword("password"),
		filemaker.SetHttpClient(http.DefaultClient),//Optional o custom
	)
	if err != nil {

	}

	serviceRecord := filemaker.NewRecordService("DatabaseName", "LayoutName", client)
	data, err := serviceRecord.GetById("1")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(data)

		serviceRecord = filemaker.NewRecordService("DatabaseName", "LayoutName", client)
		
        //List(offset, limit, sorter)
		data, err = serviceRecord.List("1", "10", filemaker.NewSorter("FieldName", filemaker.Descending))
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(data)
	

	searchService := filemaker.NewSearchService("DatabaseName", "LayoutName", client)
	data, err = searchService.
		Queries(filemaker.NewQueryFieldOperator("FieldName", "LUIS ALBERTO", filemaker.Equal)).
		SetOffset("1").SetLimit("10").
		//Sorters(filemaker.NewSorter("Sol_Nombres", filemaker.Descending)).
		Do()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(data)

}
```

## LICENSE

MIT-LICENSE. See [LICENSE](http://olivere.mit-license.org/)
or the LICENSE file provided in the repository for details.