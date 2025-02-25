# Retrieve

A lightweight and efficient Golang package for downloading files from the web with minimal code

## Examples

```go
package main

import (
    "log"
    "github.com/ciathefed/retrieve"
)

func main() {
    err := retrieve.New("https://example.com").
        SetOutput("filename.ext").
        Exec()
    if err != nil {
        log.Fatalf("failed to download file: %v", err)
    }
}
```

You can find all the examples [here](https://github.com/ciathefed/retrieve/blob/main/examples)
