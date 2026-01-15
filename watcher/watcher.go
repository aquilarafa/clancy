package watcher

import (
	"bufio"
	"io"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
)

// Watcher performs tail -f on a JSONL file
type Watcher struct {
	filePath string
	lines    chan []byte
	errors   chan error
	done     chan struct{}
}

// New creates a new file watcher
func New(filePath string) *Watcher {
	return &Watcher{
		filePath: filePath,
		lines:    make(chan []byte, 100),
		errors:   make(chan error, 1),
		done:     make(chan struct{}),
	}
}

// Lines returns the channel for new lines
func (w *Watcher) Lines() <-chan []byte {
	return w.lines
}

// Errors returns the channel for errors
func (w *Watcher) Errors() <-chan error {
	return w.errors
}

// Start begins watching the file
func (w *Watcher) Start() error {
	// Open file
	file, err := os.Open(w.filePath)
	if err != nil {
		return err
	}

	// Set up fsnotify
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		file.Close()
		return err
	}

	if err := watcher.Add(w.filePath); err != nil {
		watcher.Close()
		file.Close()
		return err
	}

	go w.watch(file, watcher)
	return nil
}

// Stop stops watching the file
func (w *Watcher) Stop() {
	close(w.done)
}

func (w *Watcher) watch(file *os.File, watcher *fsnotify.Watcher) {
	defer file.Close()
	defer watcher.Close()
	defer close(w.lines)

	reader := bufio.NewReader(file)

	// First, read all existing content
	w.readAvailable(reader)

	// Then watch for changes
	for {
		select {
		case <-w.done:
			return

		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				w.readAvailable(reader)
			}
			// Handle file truncation/recreation
			if event.Op&fsnotify.Remove == fsnotify.Remove || event.Op&fsnotify.Rename == fsnotify.Rename {
				// File was removed/renamed, try to reopen
				time.Sleep(100 * time.Millisecond)
				newFile, err := os.Open(w.filePath)
				if err != nil {
					continue
				}
				file.Close()
				file = newFile
				reader = bufio.NewReader(file)
				w.readAvailable(reader)
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			select {
			case w.errors <- err:
			default:
			}
		}
	}
}

func (w *Watcher) readAvailable(reader *bufio.Reader) {
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				// If we got a partial line at EOF, emit it (file may be complete)
				if len(line) > 0 {
					// Trim any trailing whitespace
					if line[len(line)-1] == '\r' {
						line = line[:len(line)-1]
					}
					if len(line) > 0 {
						select {
						case w.lines <- line:
						case <-w.done:
							return
						}
					}
				}
				return
			}
			select {
			case w.errors <- err:
			default:
			}
			return
		}

		// Trim newline
		if len(line) > 0 && line[len(line)-1] == '\n' {
			line = line[:len(line)-1]
		}
		if len(line) > 0 && line[len(line)-1] == '\r' {
			line = line[:len(line)-1]
		}

		if len(line) > 0 {
			select {
			case w.lines <- line:
			case <-w.done:
				return
			}
		}
	}
}
