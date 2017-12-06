package main

import (
	"flag"
	"log"
	"os"
)

func main() {
	var port = flag.Int("port", 8080, "HTTP port")
	var configFile = flag.String("config", "evil.json", "Reader config file")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.LstdFlags)

	logger.Printf("Port %d\n", *port)
	logger.Printf("Config file %s\n", *configFile)
}
