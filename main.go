package main

import (
	"flag"

	"github.com/dynamicgo/config"
	"github.com/goany/slf4go"
	"github.com/inwecrypto/neo-insight/insight"
	_ "github.com/lib/pq"
)

var logger = slf4go.Get("neo-indexer")
var configpath = flag.String("conf", "./insight.json", "neo indexer config file path")

func main() {

	flag.Parse()

	neocnf, err := config.NewFromFile(*configpath)

	if err != nil {
		logger.ErrorF("load neo config err , %s", err)
		return
	}

	server, err := insight.NewServer(neocnf)

	if err != nil {
		logger.ErrorF("load neo config err , %s", err)
		return
	}

	server.Run()
}
