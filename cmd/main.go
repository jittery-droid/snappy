package main

import (
	"context"
	"log"

	"github.com/jittery-droid/snappy/cmd/root"
)

func main() {
	if err := root.Execute(); err != nil {
		log.Fatal(context.Background(), err)
	}
}
