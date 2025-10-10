# Go-Argus

Handle file changes and reload go values based on the file contents.

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