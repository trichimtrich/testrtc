package main

import (
	"flag"
)

func main() {
	ServerHost := flag.String("server", "", "Default is 'server'. If this field is set, run as 'client'")
	ListenHost := flag.String("host", "localhost:5000", "")

	flag.Parse()

	if *ServerHost == "" {
		runServer(*ListenHost)
	} else {
		runClient(*ServerHost)
	}
}

