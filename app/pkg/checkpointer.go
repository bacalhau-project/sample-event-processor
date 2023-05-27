package app

import (
	"fmt"
	"os"
	"strconv"
)

// FileBasedCheckpointer is a Checkpointer that stores the checkpoint in a file.
type FileBasedCheckpointer struct {
	filePath string
}

func NewFileBasedCheckpointer(filePath string) *FileBasedCheckpointer {
	return &FileBasedCheckpointer{filePath: filePath}
}

func (f FileBasedCheckpointer) Restore() (int64, error) {
	// Check if the file exists
	stat, err := os.Stat(f.filePath)
	if os.IsNotExist(err) {
		return 0, nil
	} else if err != nil {
		return 0, fmt.Errorf("error checking checkpoint file: %w", err)
	} else if stat.IsDir() {
		return 0, fmt.Errorf("checkpoint file is a directory")
	}

	checkpointContent, err := os.ReadFile(f.filePath)
	if err != nil {
		return 0, fmt.Errorf("error reading checkpoint file: %w", err)
	}
	if len(checkpointContent) == 0 {
		return 0, nil
	}

	checkpoint, err := strconv.ParseInt(string(checkpointContent), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("error parsing checkpoint file: %w", err)
	}
	return checkpoint, nil
}

func (f FileBasedCheckpointer) Save(offset int64) error {
	file, err := os.OpenFile(f.filePath, os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("error creating checkpoint file: %w", err)
	}
	defer file.Close()

	// Write the int64 value to the file
	_, err = file.WriteString(strconv.FormatInt(offset, 10))
	if err != nil {
		return fmt.Errorf("error writing checkpoint file: %w", err)
	}
	return nil
}

// compile-time check for interface implementation
var _ Checkpointer = (*FileBasedCheckpointer)(nil)
