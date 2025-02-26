package main

import (
	"log"

	"github.com/wyattcupp/codesnap/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
