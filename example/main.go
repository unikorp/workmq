package main

import (
	"fmt"
	"time"

	"github.com/unikorp/workmq"
)

func main() {
	app := workmq.Init()

	app.AddProcessor("processor.logger.1s", func(worker *workmq.Worker, message workmq.Message) {
		time.Sleep(time.Second * 1)

		fmt.Printf("Worker #%d (queue: \"%s\") manages message %s\n", worker.ID, worker.Queue, message.Body)
	})

	app.AddProcessor("processor.logger.2s", func(worker *workmq.Worker, message workmq.Message) {
		time.Sleep(time.Second * 2)

		fmt.Printf("Worker #%d (queue: \"%s\") manages message %s\n", worker.ID, worker.Queue, message.Body)
	})

	app.Handle()
}
