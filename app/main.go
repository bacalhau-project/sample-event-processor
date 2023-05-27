package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/bacalhau-project/sample-event-processor/pkg"
	"github.com/hpcloud/tail"
)

func main() {
	// Parse command line arguments
	filePath := flag.String("file", "", "Path to the events file")
	checkpointPath := flag.String("checkpoint", "", "Path to the checkpoint file")
	timeout := flag.Duration("timeout", 5*time.Minute,
		"Operation timeout duration. Should be lower than Bacalhau timeout to avoid interrupting and treating the job as failed")
	flag.Parse()

	if *filePath == "" {
		log.Fatal("Please provide a valid events file path")
	}

	if *checkpointPath == "" {
		log.Fatal("Please provide a valid checkpoint file path")
	}

	if *timeout <= 0 {
		log.Fatal("Please provide a valid timeout")
	}

	processor := app.NewLoggingEventProcessor()
	checkpointer := app.NewFileBasedCheckpointer(*checkpointPath)
	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	err := run(ctx, *filePath, processor, checkpointer)
	if err != nil {
		log.Fatal(err)
	}

	// Read lines from the file until the timeout or completion
	for {
		select {
		case <-ctx.Done():
			return
		}
	}
}

func run(ctx context.Context, filePath string, processor app.EventProcessor, checkpointer app.Checkpointer) error {
	seekOffset, err := checkpointer.Restore()
	if err != nil {
		return err
	}

	// Open the event file
	t, err := tail.TailFile(filePath, tail.Config{
		Location: &tail.SeekInfo{Offset: seekOffset, Whence: io.SeekStart}, // Seek to the provided offset
		Follow:   true,                                                     // Continue tailing as new lines are appended
	})

	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}

	offset := seekOffset
	for {
		select {
		case line := <-t.Lines:
			// Received a new line from the file
			// Process each log line
			event := app.Event{
				Source: filePath,
				Offset: offset,
				Data:   line.Text,
			}
			offset += int64(len(line.Text)) + 1 // +1 for the newline character
			err = processor.ProcessEvent(event)
			if err != nil {
				return fmt.Errorf("error processing event: %w", err)
			}
			err = checkpointer.Save(offset)
			if err != nil {
				return err
			}
		case <-ctx.Done():
			// Context canceled, stop tailing the file
			t.Stop()
			return nil
		}
	}
}
