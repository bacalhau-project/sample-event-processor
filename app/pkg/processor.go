package app

import (
	"log"
)

// LoggingEventProcessor is an EventProcessor that logs each event to stdout.
type LoggingEventProcessor struct {
}

func NewLoggingEventProcessor() *LoggingEventProcessor {
	return &LoggingEventProcessor{}
}

func (p *LoggingEventProcessor) ProcessEvent(event Event) error {
	// Process each log line
	log.Printf("%+v\n", event)
	return nil
}

// compile-time check for interface implementation
var _ EventProcessor = (*LoggingEventProcessor)(nil)
