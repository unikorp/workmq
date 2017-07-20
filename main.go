package workmq

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/paulbellamy/ratecounter"
)

// RateCounters containers counters
type RateCounters map[string]*ratecounter.RateCounter

// Workmq type
type Workmq struct {
	Config     Config
	Queues     map[string]chan Message
	Processors map[string]Processor
	Counters   RateCounters
	Workers    []Worker
	Wg         sync.WaitGroup
}

// Init initializes processor part
func Init() *Workmq {
	config := GetConfig()
	processors := make(map[string]Processor)
	queues := make(map[string]chan Message)

	counters := RateCounters{
		"sent": ratecounter.NewRateCounter(1 * time.Second),
	}

	return &Workmq{
		Config:     config,
		Queues:     queues,
		Processors: processors,
		Counters:   counters,
	}
}

// Handle handles the configuration and runs workers
func (w *Workmq) Handle() {
	for queue, data := range w.Config.Queues {
		w.Queues[queue] = make(chan Message, 100000)
		w.Counters[queue] = ratecounter.NewRateCounter(1 * time.Second)

		for i := 1; i <= data.NumWorkers; i++ {
			processor, err := w.GetProcessor(data.Processor)

			if err != nil {
				panic("Unable to find processor: " + data.Processor)
			}

			current := NewWorker(i, queue, processor, w.Queues[queue], w.Counters[queue])
			w.Workers = append(w.Workers, current)

			go current.Process()
		}
	}

	w.Wg.Add(2)

	go w.ListenUDP()
	go w.ListenHTTP()

	w.Wg.Wait()
}

// ListenUDP creates a UDP server that listens for new messages
func (w *Workmq) ListenUDP() {
	defer w.Wg.Done()

	address, _ := net.ResolveUDPAddr("udp", w.Config.Ports.UDP)
	connection, _ := net.ListenUDP("udp", address)

	defer connection.Close()

	buf := make([]byte, 1024)

	for {
		n, _, _ := connection.ReadFromUDP(buf)
		w.Counters["sent"].Incr(1)

		message := TransformStringToMessage(buf[0:n])
		w.Queues[message.Queue] <- message
	}
}

// ListenHTTP creates a HTTP server to expose statistics information
func (w *Workmq) ListenHTTP() {
	defer w.Wg.Done()

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintln(writer, fmt.Sprintf("Sent rate: %d/s", w.Counters["sent"].Rate()))

		var keys []string
		for key := range w.Queues {
			keys = append(keys, key)
		}

		sort.Strings(keys)

		for _, key := range keys {
			fmt.Fprintln(writer, fmt.Sprintf("\n-> %s (%d workers):", key, w.Config.Queues[key].NumWorkers))
			fmt.Fprintln(writer, fmt.Sprintf("	Acknowledge: %d/s", w.Counters[key].Rate()))
			fmt.Fprintln(writer, fmt.Sprintf("	Messages: %d", len(w.Queues[key])))
		}
	})

	err := http.ListenAndServe(w.Config.Ports.HTTP, nil)

	if err != nil {
		log.Fatal("ListenAndServe error: ", err)
	}
}
