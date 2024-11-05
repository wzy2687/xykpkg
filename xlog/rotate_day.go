package xlog

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// pre_yymmdd.log
type writer struct {
	preName    string
	fileName   string
	mu         sync.Mutex
	MaxSize    int
	file       *os.File
	size       int64
	rotateTime time.Time
}

const megabyte = 1024 * 1024
const defaultMaxSize = 300 * 1024

func TsDayEndAt(t time.Time) time.Time {
	n := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, t.Location())
	return n
}

func newFileName(pre string) string {
	fname := pre + time.Now().Format("_20060102") + ".log"
	return fname
}

func newWriter(preName string) *writer {
	rt := &writer{preName: preName}
	rt.fileName = newFileName(rt.preName)
	rt.rotateTime = TsDayEndAt(time.Now())
	return rt
}

func (l *writer) max() int64 {
	if l.MaxSize == 0 {
		return int64(defaultMaxSize * megabyte)
	}
	return int64(l.MaxSize) * int64(megabyte)
}

// dir returns the directory for the current filename.
func (l *writer) dir() string {
	return filepath.Dir(l.fileName)
}

func (l *writer) openNew() error {
	err := os.MkdirAll(l.dir(), 0755)
	if err != nil {
		return fmt.Errorf("can't make directories for new logfile: %s", err)
	}

	name := l.fileName
	mode := os.FileMode(0600)

	// we use truncate here because this should only get called when we've moved
	// the file ourselves. if someone else creates the file in the meantime,
	// just wipe out the contents.
	f, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return fmt.Errorf("can't open new logfile: %s", err)
	}
	l.file = f
	l.size = 0
	return nil
}

func (l *writer) openExistingOrNew(writeLen int) error {

	filename := l.fileName
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return l.openNew()
	}
	if err != nil {
		return fmt.Errorf("error getting log file info: %s", err)
	}

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		// if we fail to open the old log file for some reason, just ignore
		// it and open a new log file.
		return l.openNew()
	}
	l.file = file
	l.size = info.Size()
	return nil
}

func (l *writer) Write(p []byte) (n int, err error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	writeLen := int64(len(p))
	if writeLen > l.max() {
		return 0, fmt.Errorf(
			"write length %d exceeds maximum file size %d", writeLen, l.max(),
		)
	}
	if l.file == nil {
		if err = l.openExistingOrNew(len(p)); err != nil {
			return 0, err
		}
	}

	ts := time.Now()
	if ts.After(l.rotateTime) {
		l.rotate()
	}

	n, err = l.file.Write(p)
	l.size += int64(n)

	return n, err

}
func (l *writer) rotate() error {

	if err := l.close(); err != nil {
		return err
	}
	l.fileName = newFileName(l.preName)
	if err := l.openNew(); err != nil {
		return err
	}
	l.rotateTime = TsDayEndAt(time.Now())

	return nil
}

func (l *writer) close() error {
	if l.file == nil {
		return nil
	}
	err := l.file.Close()
	l.file = nil
	return err
}
