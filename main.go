package main

import (
	"log"
	"os"

	"github.com/drj613/metrognome/cmd/metrognome"
)

func main() {
	if err := metrognome.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
