package controller

import (
	"fmt"
	"net"
	"net/rpc"
	"time"
)

const (
	PORT = ":2121"
)

type FileMetadata struct {
	Name       string
	Size       int64
	Created    time.Time
	Modified   time.Time
	ChunkAddrs []string
}

func Start(quit chan bool) {
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
