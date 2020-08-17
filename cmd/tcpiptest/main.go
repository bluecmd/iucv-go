package main

import (
	"log"

	"github.com/bluecmd/iucv-go/vmtcpip"
)

func main() {
	t, err := vmtcpip.NewTCPIP("TCPIP", "GOTEST")
	if err != nil {
		log.Fatalf("Failed to initialize VM TCPIP: %v", err)
	}
	log.Printf("VM TCPIP Hostname: %q", t.Hostname())
	log.Printf("Done")
}
