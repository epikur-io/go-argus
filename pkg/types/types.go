package types

import "io"

type IArgus interface {
	LoadValue() error
	StartWatcher() error
	StopWatcher()
}

type Decoder interface {
	Decode(val any) error
}

type ReaderFactory interface {
	NewReader() (io.ReadCloser, error)
}

type Watcher interface {
	Start(IArgus) error
	Stop()
}
