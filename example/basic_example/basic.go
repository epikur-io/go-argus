package main

import (
	"fmt"
	"log"
	"time"

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
	config, err := argus.NewArgus[Config](
		argus.WithLoader(filereader.New(configFile)),
		argus.WithWatcher(filewatcher.New(configFile)),
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
