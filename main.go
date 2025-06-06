package main

import (
	"log"
	"os"

	"github.com/djdjo/metrognome/cmd/metrognome"
)

func main() {
	if err := metrognome.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}