package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/rs/zerolog"

	argus "github.com/epikur-io/go-argus"
	filereader "github.com/epikur-io/go-argus/pkg/reader/filereader"
	filewatcher "github.com/epikur-io/go-argus/pkg/watcher/filewatcher"
)

type Config struct {
	ServerPort int    `json:"serverPort" yaml:"serverPort"`
	LogLevel   string `json:"logLevel" yaml:"logLevel"`
}

func main() {
	configFile := "./example/testfile.yaml"
	logger := zerolog.New(os.Stdin)
	config, err := argus.NewArgus[Config](
		argus.WithLoader(filereader.New(configFile)),
		argus.WithWatcher(filewatcher.New(configFile)),
		argus.WithYamlDecoder(),
		argus.WithLogger(logger),
		argus.WithCallback(func(logger *zerolog.Logger) {
			logger.Info().Msg("Successfully reloaded config!")
		}),
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

	// In your application code, you can safely access the config:
	for i := 1; i <= 5; i++ {
		go func() {
			for {
				cfg := config.GetValue()
				fmt.Printf("#%d - Server port: %d, Log level: %s\n", i,
					cfg.ServerPort, cfg.LogLevel)
				time.Sleep(3 * time.Second)
			}
		}()
	}

	// Keep main running
	select {}
}
