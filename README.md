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

## Contributing

Contributions are welcome! Please follow these steps to contribute:

1. Fork the repository.
2. Create a new branch (`git checkout -b feature/YourFeature`).
3. Make your changes and commit them (`git commit -m 'Add some feature'`).
4. Push to the branch (`git push origin feature/YourFeature`).
5. Open a pull request.

## License

This project is licensed under the MIT License - see the [LICENSE](https://github.com/ciathefed/retrieve/blob/main/LICENSE) file for details.
