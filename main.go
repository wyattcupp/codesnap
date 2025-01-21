package main

import (
	"log"

	"github.com/wyattcupp/codebase-tool/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
