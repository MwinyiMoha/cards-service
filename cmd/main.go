package main

import (
	"log"

	"github.com/mwinyimoha/commons/pkg/logging"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	logger, err := logging.NewLoggerConfig().BuildLogger()
	if err != nil {
		log.Fatal("could not initialize logging", err)
	}

	defer logger.Sync()
}
