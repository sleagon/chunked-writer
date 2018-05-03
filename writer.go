package writer

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sync"
	"time"
)

const (
	// DefaultChunkSize default chunk size 16M
	DefaultChunkSize = 1 << 24
	// DefaultPrefix default prefix
	DefaultPrefix = "default"
	// DefaultPattern default pattern for date
	DefaultPattern = "20060102"
	// DefaultDir default dir current directory
	DefaultDir = "."
)

// Writer chunked writer
type Writer struct {
	sync.RWMutex
	ChunkSize int64
	Prefix    string
	Pattern   string
	Dir       string
	no        int64
}

// New get a new writer
func New(dir string, prefix string, size int64) (*Writer, error) {
	w := new(Writer)
	w.Dir = dir
	w.Prefix = prefix
	w.ChunkSize = size
	w.Pattern = "20060102"
	w.Default()
	if err := os.MkdirAll(w.Dir, 0700); err != nil {
		return nil, err
	}
	return w, nil
}

// Default set writer default value
func (w *Writer) Default() {
	if w.ChunkSize == 0 {
		w.ChunkSize = DefaultChunkSize // default
	}
	if w.Prefix == "" {
		w.Prefix = DefaultPrefix
	}
	if w.Pattern == "" {
		w.Pattern = DefaultPattern
	}
	if w.Dir == "" {
		w.Dir = DefaultDir
	}
}

// Add add no of w
func (w *Writer) add(n int64) int64 {
	w.Lock()
	defer w.Unlock()
	w.no += n
	return w.no
}

func (w *Writer) reset() {
	w.Lock()
	defer w.Unlock()
	w.no = 0
}

func (w *Writer) pickName(fp string) (string, error) {
	for {
		no := w.add(1)
		n := fmt.Sprintf("%s.%d", fp, no)
		nex, err := isNotExist(n)
		if err != nil {
			return "", err
		}
		if nex {
			return n, nil
		}
	}
}

func (w *Writer) open(fp string, reserved int64) (*os.File, error) {
	nex, err := isNotExist(fp)
	if err != nil {
		return nil, err
	}
	if nex {
		return os.Create(fp)
	}
	stat, err := os.Stat(fp)
	if err != nil {
		return nil, err
	}
	if stat.Size() <= w.ChunkSize-reserved {
		return os.OpenFile(fp, os.O_APPEND|os.O_WRONLY, 0)
	}
	n, err := w.pickName(fp)
	if err := os.Rename(fp, n); err != nil {
		w.reset()
		return nil, err
	}
	return os.Create(fp)
}

// Write realize io.Writer
func (w *Writer) Write(p []byte) (int, error) {
	day := time.Now().Format(w.Pattern)
	file := fmt.Sprintf("%s.%s.log", w.Prefix, day)
	size := int64(len(p))
	fp := path.Join(w.Dir, file)
	fw, err := w.open(fp, size)
	if err != nil {
		return 0, err
	}
	return fw.Write(p)
}

func isNotExist(fp string) (bool, error) {
	_, err := os.Stat(fp)
	if err == nil {
		return false, nil
	}
	if os.IsNotExist(err) {
		return true, nil
	}
	return false, err
}

func removeAll(pf string) error {
	files, err := filepath.Glob(pf)
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		if err := os.Remove(f); err != nil {
			return err
		}
	}
	return nil
}
