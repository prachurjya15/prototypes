package main

import (
	"sync"

	"github.com/prachurjya15/prototypes/distributed-id-generators/snowflake/services"
)

func main() {
	// Suppose we have three services, modeled as goroutines. Each service generates Ids.
	svcs := []services.Service{
		*services.NewService("Service-A"),
		*services.NewService("Service-B"),
	}

	var wg sync.WaitGroup
	for _, s := range svcs {
		wg.Go(s.Work)
	}
	wg.Wait()
}
