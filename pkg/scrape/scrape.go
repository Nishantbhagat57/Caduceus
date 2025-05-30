package scrape

import (
	"fmt"
	"net"

	"time"

	"github.com/Nishantbhagat57/Caduceus/pkg/types"
	"github.com/Nishantbhagat57/Caduceus/pkg/utils"
	"github.com/Nishantbhagat57/Caduceus/pkg/workers"
)

func RunScrape(args types.ScrapeArgs) {
	dialer := &net.Dialer{
		Timeout: time.Duration(args.Timeout) * time.Second,
	}

	inputChannel := make(chan string)
	resultChannel := make(chan types.Result)
	outputChannel := make(chan string, args.Concurrency/10)

	// Create and start the worker pool
	workerPool := workers.NewWorkerPool(args.Concurrency, dialer, inputChannel, resultChannel)
	workerPool.Start()

	// Create and start the results worker pool
	resultsWorkerPool := workers.NewResultWorkerPool(args.Concurrency/100, resultChannel, outputChannel) // Adjust the size as needed
	resultsWorkerPool.Start(args)

	// Handle input feeding
	go func() {
		utils.IntakeFunction(inputChannel, args.Ports, args.Input)
		close(inputChannel)
	}()

	// Handle outputs
	go func() {
		for output := range outputChannel {
			fmt.Println(output)
		}
	}()

	workerPool.Stop()
	resultsWorkerPool.Stop()

	// if args.PrintStats {
	// 	stats.Display() // Display updated stats
	// }
}
