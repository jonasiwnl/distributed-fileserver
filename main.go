package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/jonasiwnl/distributed-fileserver/v2/server"
)

const (
	CONTROLLERPORT = ":2120"
	FILESERVERPORT = ":2125"
)

func main() {
	flagController := flag.Bool("controller", false, "Start controller server.")
	flagFileServer := flag.Bool("fileserver", false, "Start fileserver.")

	flag.Parse()

	interrupt := make(chan os.Signal, 1)
	quit := make(chan bool, 1)

	// Quit fileserver on interrupt
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-interrupt
		quit <- true
	}()

	if *flagController {
		server.StartControllerServer(CONTROLLERPORT, quit)
	} else if *flagFileServer {
		sixtyFourMB := int64(64 * 1024 * 1024)
		server.StartFileServer(FILESERVERPORT, sixtyFourMB, quit)
	}
}
