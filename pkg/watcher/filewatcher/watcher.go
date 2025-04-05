package filewatcher

import (
	"github.com/epikur-io/go-argus/pkg/types"

	"github.com/fsnotify/fsnotify"
	"github.com/rs/zerolog"
)

type option func(*Watcher)

func WithLogger(l zerolog.Logger) option {
	return func(c *Watcher) {
		c.logger = l
		c.hasLogger = true
	}
}

func New(file string, opts ...option) *Watcher {
	w := &Watcher{
		file:        file,
		stopWatcher: make(chan struct{}),
	}
	// apply options
	for _, opt := range opts {
		if opt != nil {
			opt(w)
		}
	}
	return w
}

type Watcher struct {
	file        string
	logger      zerolog.Logger
	hasLogger   bool
	stopWatcher chan struct{}
}

func (w *Watcher) Start(argus types.IArgus) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	watcher.Add(w.file)
	go func() {
		defer watcher.Close()
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if w.hasLogger {
					w.logger.Debug().Str("event", event.Name).Str("Op", event.Op.String()).Msg("got event from file watcher")
				}
				if event.Has(fsnotify.Write) {
					if err := argus.LoadValue(); err != nil {
						if w.hasLogger {
							w.logger.Error().Err(err).Msg("failed to load value")
						}
						continue
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				if w.hasLogger {
					w.logger.Error().Err(err).Msg("HotValue file watcher error")
				}
			case <-w.stopWatcher:
				if w.hasLogger {
					w.logger.Debug().Msg("stopped file watcher")
				}
				return
			}
		}
	}()
	return nil
}

func (w *Watcher) Stop() {
	w.stopWatcher <- struct{}{}
}
