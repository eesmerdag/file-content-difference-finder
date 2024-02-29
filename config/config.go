package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	Port        int
	FileContent string
	FileVersion int
}

func InitConfig() *Config {
	config := Config{}

	port, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalf("invalid port parameter, Err: %s", err)
	}
	config.Port = port

	fileVersion, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatalf("invalid file version, Err: %s", err)
	}
	config.FileVersion = fileVersion
	config.FileContent = os.Args[3]

	return &config
}
