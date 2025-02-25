package main

import (
	"log"

	"github.com/ciathefed/retrieve"
)

func main() {
	err := retrieve.New("https://cataas.com/cat").
		SetOutput("cat.png").
		Exec()
	if err != nil {
		log.Fatalf("failed to download file: %v", err)
	}
}
