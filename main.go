package main

import (
	"flag"

	log "github.com/sirupsen/logrus"
)

func main() {

	log.SetFormatter(&log.TextFormatter{})

	// parse command line args
	configFilePath := flag.String("c", Canonicalize("~/.tagger.yml"), "Path to config file")
	debug := flag.Bool("d", false, "Activate debug logging level")
	x := flag.Bool("x", false, "Generate tags in current directory and exit")
	flag.Parse()

	if *debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	// parse config
	var config *Config
	if *x {
		config = &Config{
			Indexer:  DefaultIndexer(),
			Projects: []struct{ Path string }{{"."}},
		}
	} else {
		config = NewConfig(*configFilePath)
	}

	manager := NewManager(config)
	go manager.Listen(config.Port)
	manager.Start()
}
