package goargus

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"sync/atomic"

	"github.com/epikur-io/go-argus/pkg/types"

	"github.com/BurntSushi/toml"
	"github.com/rs/zerolog"
	"gopkg.in/ini.v1"
	yaml "gopkg.in/yaml.v3"
)

type config struct {
	logger        zerolog.Logger
	hasLogger     bool
	watcher       types.Watcher
	readerFactory types.ReaderFactory
	callback      func(l *zerolog.Logger)
	decoder       func(r io.Reader) types.Decoder
}

type option func(*config)

func WithLogger(l zerolog.Logger) option {
	return func(c *config) {
		c.logger = l
		c.hasLogger = true
	}
}

// Callback fn is executed on sucessful reload of the value
func WithCallback(fn func(l *zerolog.Logger)) option {
	return func(c *config) {
		c.callback = fn
	}
}

// file decoders
func WithJsonDecoder() option {
	return func(c *config) {
		c.decoder = func(r io.Reader) types.Decoder {
			return json.NewDecoder(r)
		}
	}
}

func WithYamlDecoder() option {
	return func(c *config) {
		c.decoder = func(r io.Reader) types.Decoder {
			return yaml.NewDecoder(r)
		}
	}
}

type iniDecoderWrapper struct {
	reader io.Reader
}

func (d *iniDecoderWrapper) Decode(val any) error {
	iniFile, err := ini.Load(d.reader)
	if err != nil {
		return err
	}
	return iniFile.MapTo(val)
}

func WithIniDecoder() option {
	return func(c *config) {
		c.decoder = func(r io.Reader) types.Decoder {
			return &iniDecoderWrapper{
				reader: r,
			}
		}
	}
}

type tomlDecoderWrapper struct {
	decoder *toml.Decoder
}

func (d *tomlDecoderWrapper) Decode(val any) error {
	_, err := d.decoder.Decode(val)
	return err
}

func WithTomlDecoder() option {
	return func(c *config) {
		c.decoder = func(r io.Reader) types.Decoder {
			return &tomlDecoderWrapper{
				decoder: toml.NewDecoder(r),
			}
		}
	}
}

// Use custom file decoder, see the Decoder interface
func WithCustomDecoder(decoderFactory func(io.Reader) types.Decoder) option {
	return func(c *config) {
		c.decoder = decoderFactory
	}
}

func WithReader(loader types.ReaderFactory) option {
	return func(c *config) {
		c.readerFactory = loader
	}
}

func WithWatcher(watcher types.Watcher) option {
	return func(c *config) {
		c.watcher = watcher
	}
}

// Argus handles hot reloading of config files in different formats
type Argus[T any] struct {
	value  atomic.Value
	mux    sync.Mutex
	config config
}

// Creates a new Argus file watcher and decode the file in to type T
func NewArgus[T any](opts ...option) (*Argus[T], error) {
	m := &Argus[T]{
		config: config{},
	}

	// apply options
	for _, opt := range opts {
		if opt != nil {
			opt(&m.config)
		}
	}
	// set yaml default decoder
	if m.config.decoder == nil {
		WithYamlDecoder()(&m.config)
	}
	if m.config.readerFactory == nil {
		return nil, fmt.Errorf("missing loader")
	}

	if err := m.LoadValue(); err != nil {
		return nil, fmt.Errorf("failed to load initial value: %s", err)
	}

	return m, nil
}

func (m *Argus[T]) InjectIntoContext(ctx context.Context, key any) context.Context {
	val := m.GetValue()
	// nolint:staticcheck
	return context.WithValue(ctx, key, val)
}

func (m *Argus[T]) GetValue() T {
	return m.value.Load().(T)
}

func (m *Argus[T]) LoadValue() error {
	m.mux.Lock()
	defer m.mux.Unlock()

	reader, err := m.config.readerFactory.NewReader()
	if err != nil {
		return err
	}
	defer reader.Close()

	var newValue T
	decoder := m.config.decoder(reader)
	if err := decoder.Decode(&newValue); err != nil {
		return err
	}

	// store the new config atomically
	m.value.Store(newValue)

	if m.config.callback != nil {
		var logger *zerolog.Logger
		if m.config.hasLogger {
			logger = &m.config.logger
		}
		m.config.callback(logger)
	}

	return nil
}

func (m *Argus[T]) StartWatcher() error {
	if m.config.watcher == nil {
		return fmt.Errorf("no watcher attached to argus")
	}
	return m.config.watcher.Start(m)
}

func (m *Argus[T]) StopWatcher() {
	m.config.watcher.Stop()
}
