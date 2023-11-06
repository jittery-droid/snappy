package main

import (
	"context"
	"log"

	"github.com/jittery-droid/snappy/cmd/root"
)

// curl http://localhost:6060/debug/pprof/heap > heap_staging.out
// http://localhost:6060/debug/pprof/profile?seconds=30
// http://localhost:6060/debug/pprof/block
// http://localhost:6060/debug/pprof/mutex
// http://localhost:6060/debug/pprof/trace?seconds=5
func main() {
	if err := root.Execute(); err != nil {
		log.Fatal(context.Background(), err)
	}
}
