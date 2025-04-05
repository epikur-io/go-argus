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

type Loader interface {
	Load() (io.ReadCloser, error)
}

type Watcher interface {
	Start(IArgus) error
	Stop()
}
