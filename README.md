# FileMaker

**This is a library created in go to connect to the FileMaker Data API.**

[FileMaker API](https://fmhelp.filemaker.com/docs/18/es/dataapi/) documentation.

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
The client connects to FileMaker Server 18.1 passing a host.


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
		filemaker.SetHttpClient(http.DefaultClient),//Optional or custom
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
		GroupQuery(filemaker.NewQueryGroup(
                   				filemaker.NewQueryFieldOperator("name", "pablo", filemaker.Equal),
                   				filemaker.NewQueryFieldOperator("last_name", "zenteno", filemaker.Equal),
                   			),
        ).Sorters(filemaker.NewSorter("last_name", filemaker.Descending)).Do()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(data)

}
```

#### Author
* **[Pablo Zenteno](https://github.com/pzentenoe)** - *Full Stack developer*

## LICENSE

MIT-LICENSE. See [LICENSE](http://olivere.mit-license.org/)
or the LICENSE file provided in the repository for details.
