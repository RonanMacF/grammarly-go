package flags

import (
	"errors"
	"flag"
)

type flags struct {
	ConfigPath string
	FilePath string
	LogPath string
	LogLevel int
}
var (
	configPath = flag.String("configPath",  "/Users/ronan/go/src/grammarly-go/sampleConfig.toml" ,
		"path to the directory containing the configuration file")
	filePath = flag.String("filePath", "/Users/ronan/tmp.txt",  "path to the file being checked")
	logPath = flag.String("logPath", "/Users/ronan/grammarly-go/grammarly-go.log",  "path to the file to print logs")
	logLevel = flag.Int("logLevel", 4,  "log level")
)

func Parse() (*flags, error) {
	flag.Parse()
	if *filePath == ""{
		return nil, errors.New("no file passed, pass file using -filePath flag")
	}
	return &flags{
		ConfigPath: *configPath,
		FilePath:   *filePath,
		LogPath: *logPath,
		LogLevel: *logLevel,
	}, nil
}

