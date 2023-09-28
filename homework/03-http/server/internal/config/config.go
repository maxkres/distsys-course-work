package config

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
	Host             string
	Port             string
	WorkingDirectory string
	Domain           string
}

func New() (Config, error) {
	var host, port, workingDirectory, serverDomain string
	flag.StringVar(&host, "host", "0.0.0.0", "server host")
	flag.StringVar(&port, "port", "8080", "server port")
	flag.StringVar(&workingDirectory, "working-directory", "", "directory for files")
	flag.StringVar(&serverDomain, "server-domain", "", "server domain")

	flag.Parse()

	if !isFlagPassed("host") {
		ex, val := envExists("SERVER_HOST")
		if ex {
			host = val
		}
	}
	if !isFlagPassed("port") {
		ex, val := envExists("SERVER_PORT")
		if ex {
			port = val
		}
	}
	if !isFlagPassed("working-directory") {
		ex, val := envExists("SERVER_WORKING_DIRECTORY")
		if ex {
			workingDirectory = val
		} else {
			return Config{}, fmt.Errorf("server working directory is required!")
		}
	}
	if !isFlagPassed("server-domain") {
		ex, val := envExists("SERVER_DOMAIN")
		if ex {
			serverDomain = val
		}
	}

	return Config{
		Host:             host,
		Port:             port,
		WorkingDirectory: workingDirectory,
		Domain:           serverDomain,
	}, nil
}

func envExists(name string) (bool, string) {
	val := os.Getenv(name)
	if val == "" {
		return false, ""
	}
	return true, val
}

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}
