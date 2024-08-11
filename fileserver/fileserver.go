package fileserver

import (
	"fmt"
	"net"
	"net/rpc"
)

const (
	PORT      = ":2122"
	DIRECTORY = "virtual/"
)

// Used by RPC handlers in fileops.go
type FileServer struct{}

func Start(quit chan bool) {
	fileServer := new(FileServer)
	rpc.Register(fileServer)

	listener, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Printf("Couldn't listen on port %s: %s\n", PORT, err)
		return
	}

	fmt.Printf("Listening on port %s\n", PORT)

	for {
		select {
		case <-quit:
			listener.Close()
			return
		default:
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("Error accepting connection: ", err)
				break
			}

			go rpc.ServeConn(conn)
		}
	}
}
