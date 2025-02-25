# Retrieve

![Test](https://img.shields.io/github/actions/workflow/status/ciathefed/retrieve/test.yml?label=test%20%F0%9F%A7%AA&style=flat-square)

A lightweight and efficient Golang package for downloading files from the web with minimal code

## Install

```shell
go get -u github.com/ciathefed/retrieve
```

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

You can find all the examples [here](https://github.com/ciathefed/retrieve/blob/main/_examples)
