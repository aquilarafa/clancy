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
	// Verify file exists
	if _, err := os.Stat(w.filePath); err != nil {
		return err
	}

	go w.watch()
	return nil
}

// Stop stops watching the file
func (w *Watcher) Stop() {
	close(w.done)
}

func (w *Watcher) watch() {
	defer close(w.lines)

	var offset int64 = 0

	for {
		// Try to watch with fsnotify
		var shouldRetry bool
		shouldRetry, offset = w.watchWithFsnotify(offset)
		if shouldRetry {
			// Session ended - wait for new session
			if !w.waitForNewSession() {
				return // done was closed
			}
		} else {
			return // done was closed
		}
	}
}

// watchWithFsnotify watches file using fsnotify. Returns (shouldRetry, newOffset).
func (w *Watcher) watchWithFsnotify(offset int64) (bool, int64) {
	file, err := os.Open(w.filePath)
	if err != nil {
		return true, offset
	}
	defer file.Close()

	// Check if file was truncated (new session with fresh file)
	info, err := file.Stat()
	if err != nil {
		return true, offset
	}
	if info.Size() < offset {
		offset = 0 // File truncated, start from beginning
	}

	// Seek to offset
	if _, err := file.Seek(offset, io.SeekStart); err != nil {
		return true, offset
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return true, offset
	}
	defer watcher.Close()

	if err := watcher.Add(w.filePath); err != nil {
		return true, offset
	}

	reader := bufio.NewReader(file)

	// Read available content from offset
	offset = w.readAvailable(reader, offset)

	// Idle timeout to detect session end
	idleTimeout := time.NewTimer(3 * time.Second)
	defer idleTimeout.Stop()

	for {
		select {
		case <-w.done:
			return false, offset

		case <-idleTimeout.C:
			// No activity for 3 seconds, session likely ended
			return true, offset

		case event, ok := <-watcher.Events:
			if !ok {
				return true, offset
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				offset = w.readAvailable(reader, offset)
				// Reset idle timeout on activity
				if !idleTimeout.Stop() {
					select {
					case <-idleTimeout.C:
					default:
					}
				}
				idleTimeout.Reset(3 * time.Second)
			}
			if event.Op&fsnotify.Remove == fsnotify.Remove || event.Op&fsnotify.Rename == fsnotify.Rename {
				return true, offset
			}

		case _, ok := <-watcher.Errors:
			if !ok {
				return true, offset
			}
		}
	}
}

// waitForNewSession polls until file changes (new content or new file). Returns false if done.
func (w *Watcher) waitForNewSession() bool {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	// Get current file info
	var lastSize int64 = -1
	var lastMod time.Time
	if info, err := os.Stat(w.filePath); err == nil {
		lastSize = info.Size()
		lastMod = info.ModTime()
	}

	for {
		select {
		case <-w.done:
			return false
		case <-ticker.C:
			info, err := os.Stat(w.filePath)
			if err != nil {
				// File removed, wait for it to reappear
				lastSize = -1
				continue
			}
			// Check if file changed (new session started)
			if info.Size() != lastSize || info.ModTime() != lastMod {
				return true
			}
		}
	}
}

func (w *Watcher) readAvailable(reader *bufio.Reader, offset int64) int64 {
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				// Partial line - don't emit, wait for newline
				return offset
			}
			return offset
		}

		// Update offset with bytes read
		offset += int64(len(line))

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
				return offset
			}
		}
	}
}
