package workmq

import (
	"fmt"

	"github.com/paulbellamy/ratecounter"
)

// Worker struct
type Worker struct {
	ID        int
	Queue     string
	Message   <-chan Message
	Processor Processor
	Counter   *ratecounter.RateCounter
}

// NewWorker creates a new Worker instance
func NewWorker(id int, queue string, processor Processor, message <-chan Message, counter *ratecounter.RateCounter) Worker {
	return Worker{ID: id, Queue: queue, Processor: processor, Message: message, Counter: counter}
}

// Process listens for a processor on the worker.
func (w *Worker) Process() {
	fmt.Printf("-> Worker %d ready to process queue \"%s\"...\n", w.ID, w.Queue)

	for message := range w.Message {
		w.Counter.Incr(1)
		w.Processor(w, message)
	}
}
