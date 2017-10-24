package main

import (
	"fmt"
	"time"

	"github.com/unikorp/workmq"
)

func handleProcessor(worker *workmq.Worker, message workmq.Message) {
	fmt.Printf("Worker #%d (queue: \"%s\") manages message %s\n", worker.ID, worker.Queue, message.Bod)
}

func main() {
	app := workmq.Init()

	app.AddProcessor("processor.logger.1s", func(worker *workmq.Worker, message workmq.Message) {
		time.Sleep(time.Second * 1)

		handleProcessor(worker, message)
	})

	app.AddProcessor("processor.logger.2s", func(worker *workmq.Worker, message workmq.Message) {
		time.Sleep(time.Second * 2)

		handleProcessor(worker, message)
	})

	app.Handle()
}
