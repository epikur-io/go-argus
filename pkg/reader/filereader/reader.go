package filereader

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

func New(file string) *Loader {
	absFilepath, err := filepath.Abs(file)
	if err == nil {
		file = absFilepath
	}
	l := &Loader{
		file: file,
	}
	return l
}

type Loader struct {
	file        string
	fh          *os.File
	lastModTime time.Time
}

func (l *Loader) NewReader() (io.ReadCloser, error) {
	var err error
	l.fh, err = os.Open(l.file)
	if err != nil {
		return nil, err
	}

	// get file info to check modification time
	fileInfo, err := l.fh.Stat()
	if err != nil {
		return nil, err
	}

	// only reload if file has changed (sanity check)
	if !fileInfo.ModTime().After(l.lastModTime) {
		return nil, fmt.Errorf("file not changed according to mtime")
	}
	l.lastModTime = fileInfo.ModTime()
	return l.fh, nil
}

func (l *Loader) Close() error {
	return l.fh.Close()
}
