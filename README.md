# Go-Argus

**Go-Argus** is a lightweight and extensible library for **reading, watching, and binding configuration data** from various sources and formats directly into Go structs.
It supports file-based configurations out of the box, with an easy way to extend to other sources (like environment variables, APIs, etc.).

## Features

- Read configuration from files (JSON, YAML, etc.)
- Watch for configuration changes and automatically reload
- Bind configuration data to typed Go structs
- Pluggable readers, watchers, and decoders
- Simple and idiomatic API

## Example

```go
package main

import (
	"fmt"
	"log"
	"time"

	argus "github.com/epikur-io/go-argus"
	"github.com/epikur-io/go-argus/pkg/reader/filereader"
	"github.com/epikur-io/go-argus/pkg/watcher/filewatcher"
)

type Config struct {
	ServerPort int    `json:"serverPort" yaml:"serverPort"`
	LogLevel   string `json:"logLevel" yaml:"logLevel"`
}

func main() {
	configFile := "./example/testfile.yaml"
	config, err := argus.NewArgus[Config](
		argus.WithReader(filereader.New(configFile)), 	// (mandatory) will return error if no reader is provided
		argus.WithWatcher(filewatcher.New(configFile)), // (optional)
		argus.WithYamlDecoder(), 						// (optional) Use YAML decoder (default)
	)
	if err != nil {
		panic(err)
	}

	// Start watching for file updates
	err = config.StartWatcher()
	if err != nil {
		log.Fatalln(err)
	}

	defer config.StopWatcher()

	for {
		cfg := config.GetValue()
		fmt.Printf("Server port: %d, Log level: %s\n", cfg.ServerPort, cfg.LogLevel)
		time.Sleep(3 * time.Second)
	}
}
```

## Run tests

Using [taskfile](https://taskfile.dev/) (see `./Taskfile.yaml` for all commands):

```bash
task test

# for test coverage:
tesk test:coverage
```
Using plain Go:

```bash
go test -v ./...
```

## Extending Go-Argus

You can easily implement your own:

- Reader: for reading config from a new source
- Watcher: to detect when the source changes
- Decoder: to handle a new data format

Each is implemented through simple interfaces that integrate seamlessly with argus.NewArgus.