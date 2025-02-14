package filewriter

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"
)

type FileWriter struct {
	file       *os.File
	writer     *bufio.Writer
	mu         sync.Mutex
	writeCount atomic.Int64
	flushCount atomic.Int64

	// Configuration
	bufferSize    int
	maxBufferSize int
	flushInterval time.Duration
	filePath      string

	// Shutdown
	done chan struct{}
}

func NewFileWriter(filePath string, options ...Option) (*FileWriter, error) {
	fw := &FileWriter{
		bufferSize:    4096,        // 4KB default
		maxBufferSize: 3072,        // 75% of buffer
		flushInterval: time.Second, // 1 second default
		filePath:      filePath,
		done:          make(chan struct{}),
	}

	// Apply options
	for _, opt := range options {
		opt(fw)
	}

	// Ensure directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("create directory failed: %v", err)
	}

	if err := fw.openFile(); err != nil {
		return nil, err
	}

	// Start background flush
	go fw.periodicFlush()

	return fw, nil
}

func (fw *FileWriter) openFile() error {
	file, err := os.OpenFile(fw.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("open file failed: %v", err)
	}

	fw.file = file
	fw.writer = bufio.NewWriterSize(file, fw.bufferSize)
	return nil
}

func (fw *FileWriter) WriteLine(line string) error {
	fw.mu.Lock()
	defer fw.mu.Unlock()

	// Write line with space separator and newline
	_, err := fmt.Fprintf(fw.writer, "%s\n", line)
	if err != nil {
		return fmt.Errorf("write failed: %v", err)
	}

	fw.writeCount.Add(1)

	// Flush if buffer is getting full
	if fw.writer.Buffered() > fw.maxBufferSize {
		if err := fw.writer.Flush(); err != nil {
			return fmt.Errorf("flush failed: %v", err)
		}
		fw.flushCount.Add(1)
	}

	return nil
}

func (fw *FileWriter) periodicFlush() {
	ticker := time.NewTicker(fw.flushInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := fw.Flush(); err != nil {
				slog.Info(err.Error())
			}
		case <-fw.done:
			return
		}
	}
}

func (fw *FileWriter) Flush() error {
	fw.mu.Lock()
	defer fw.mu.Unlock()

	if fw.writer.Buffered() > 0 {
		if err := fw.writer.Flush(); err != nil {
			return fmt.Errorf("flush failed: %v", err)
		}
		fw.flushCount.Add(1)
	}
	return nil
}

func (fw *FileWriter) Close() error {
	// Signal periodic flush to stop
	close(fw.done)

	fw.mu.Lock()
	defer fw.mu.Unlock()

	// Flush remaining data
	if err := fw.writer.Flush(); err != nil {
		return fmt.Errorf("final flush failed: %v", err)
	}

	// Close file
	if err := fw.file.Close(); err != nil {
		return fmt.Errorf("file close failed: %v", err)
	}

	return nil
}

// GetStats returns current statistics
func (fw *FileWriter) GetStats() (writeCount, flushCount int64) {
	return fw.writeCount.Load(), fw.flushCount.Load()
}

// Option pattern for configuration
type Option func(*FileWriter)

func WithBufferSize(size int) Option {
	return func(fw *FileWriter) {
		fw.bufferSize = size
		fw.maxBufferSize = size * 3 / 4
	}
}

func WithFlushInterval(interval time.Duration) Option {
	return func(fw *FileWriter) {
		fw.flushInterval = interval
	}
}
