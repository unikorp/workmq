package workmq

import "fmt"

// Processor type
type Processor func(worker *Worker, message Message)

// AddProcessor adds a processor into the processors list
func (w *Workmq) AddProcessor(name string, processor Processor) {
	w.Processors[name] = processor
}

// GetProcessor retrieves a processor from its name
func (w *Workmq) GetProcessor(name string) (Processor, error) {
	if _, ok := w.Processors[name]; !ok {
		return nil, fmt.Errorf("Unable to find processor '%s'", name)
	}

	return w.Processors[name], nil
}

// RemoveProcessor removes a processor from its name
func (w *Workmq) RemoveProcessor(name string) {
	delete(w.Processors, name)
}
