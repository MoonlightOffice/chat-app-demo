package config

import (
	"fmt"
	"log"
	"os"
)

func init() {
	err := LoadConfig()
	if err != nil {
		msg := fmt.Sprintf("failed to load config: %s", err.Error())
		log.Fatal(msg)
		os.Exit(1)
	}
}
