package main

import (
	"fmt"
	"godis/config"
	"godis/lib/logger"
	"godis/tcp"
	"os"
)

// configFile is the name of the config file
const configFile string = "godis.conf"

// set default config to use if no config file is found
var defaultProperties = &config.ServerProperties{
	Bind: "0.0.0.0",
	Port: 6379,
}

// fileExists checks if a file exists and is not a directory before we try
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	return err == nil && !info.IsDir()
}

func main() {
	// setup logger
	logger.Setup(&logger.Settings{
		Path:       "logs",
		Name:       "godis",
		Ext:        "log",
		TimeFormat: "2006-01-02",
	})

	// setup config
	if fileExists(configFile) {
		config.SetupConfig(configFile)
	} else {
		config.Properties = defaultProperties
	}

	// start server
	err := tcp.ListenAndServeWithSignal(
		&tcp.Config{
			Address: fmt.Sprintf("%s:%d", config.Properties.Bind, config.Properties.Port),
		},
		tcp.NewHandler(),
	)
	if err != nil {
		logger.Error(err)
	}

}
