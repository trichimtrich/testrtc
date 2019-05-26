package main

import (
	"flag"
)

func main() {
	ServerHost := flag.String("server", "", "Mode of operation. Default is 'SERVER'. Otherwise as 'CLIENT', ex: 'localhost:5000', 'ws://localhost:5000', ...")
	PartnerID := flag.String("partner", "", "For 'CLIENT' mode: ID of partner.")
	ListenHost := flag.String("host", "localhost:5000", "For 'SERVER' mode: Interface and port")

	flag.Parse()

	if *ServerHost == "" {
		runServer(*ListenHost)
	} else {
		runClient(*ServerHost, *PartnerID)
	}
}

