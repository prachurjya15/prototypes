package main

import (
	"sync"

	"github.com/prachurjya15/prototypes/distributed-id-generators/disk-id-store/services"
)

func main() {
	// Main invokes N services (go-routines in my case)
	// Suppose Each go-routine is a API Server asking THAT needs unique monotonous Id for some insert logic

	//Load the
	connStr := "localhost:8085"
	svc := []services.Service{
		*services.NewService("Service-A", connStr),
		*services.NewService("Service-B", connStr),
		*services.NewService("Service-C", connStr),
		*services.NewService("Service-D", connStr),
		*services.NewService("Service-E", connStr),
		*services.NewService("Service-F", connStr),
		*services.NewService("Service-G", connStr),
	}

	var wg sync.WaitGroup
	for _, s := range svc {
		wg.Go(s.Work)
	}
	wg.Wait()
}
