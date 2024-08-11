package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/jonasiwnl/distributed-fileserver/v2/fileserver"
)

func main() {
	interrupt := make(chan os.Signal, 1)
	quit := make(chan bool, 1)

	// Quit fileserver on interrupt
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-interrupt
		quit <- true
	}()

	fileserver.Start(quit)
}
