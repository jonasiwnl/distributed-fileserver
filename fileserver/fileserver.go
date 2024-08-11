package fileserver

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
	"net/rpc"
)

const (
	PORT      = ":2122"
	CTRLPORT  = ":9999"
	DIRECTORY = "virtual/"
)

// Used by RPC handlers in fileops.go
type FileServer struct {
	SizeUsed int64
	Capacity int64
}

func Start(capacity int64, quit chan bool) {
	fileServer := new(FileServer)
	rpc.Register(fileServer)

	listener, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Printf("Couldn't listen on port %s: %s\n", PORT, err)
		return
	}
	defer listener.Close()

	fmt.Printf("Fileserver listening on port %s\n", PORT)

	// Let controller know we're here
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err = encoder.Encode(FileServer{0, capacity})
	if err != nil {
		fmt.Println("Error encoding data: ", err)
		return
	}

	conn, err := net.Dial("udp4", "255.255.255.255"+CTRLPORT)
	if err != nil {
		fmt.Println("Error connecting to controller: ", err)
		return
	}

	_, err = conn.Write([]byte(buffer.Bytes()))
	if err != nil {
		fmt.Println("Error writing to controller: ", err)
		return
	}
	conn.Close()

	// Finally, just listen for file operations
	for {
		select {
		case <-quit:
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
