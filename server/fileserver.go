package server

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
	"net/rpc"
)

const (
	CTRLJOINPORT = ":9998"
	DIRECTORY    = "virtual/"
)

type FileServer struct{}

func StartFileServer(port string, capacity int64, quit chan bool) {
	fileServer := new(FileServer)
	rpc.Register(fileServer)

	listener, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Printf("Couldn't listen on port %s: %s\n", port, err)
		return
	}
	defer listener.Close()

	fmt.Printf("Fileserver listening on port %s\n", port)

	// Let controller know we're here
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	message := FileServerMessage{REGISTER, FileServerData{0, capacity}}
	err = encoder.Encode(message)
	if err != nil {
		fmt.Println("Error encoding data: ", err)
		return
	}

	conn, err := net.Dial("udp4", "255.255.255.255"+CTRLJOINPORT)
	if err != nil {
		fmt.Println("Error connecting to controller: ", err)
		return
	}
	defer conn.Close()

	_, err = conn.Write([]byte(buffer.Bytes()))
	if err != nil {
		fmt.Println("Error writing register message: ", err)
		return
	}
	// TODO: listen for ACK?

	// This goroutine cleans everything up when we hear the quit signal.
	go func() {
		// Wait for quit signal
		<-quit

		// Let controller know we're leaving
		buffer.Reset()
		encoder = gob.NewEncoder(&buffer)
		message = FileServerMessage{DISCONNECT, FileServerData{0, 0}}
		err = encoder.Encode(message)
		if err != nil {
			fmt.Println("Error encoding disconnect message: ", err)
			return
		}
		conn.Write([]byte(buffer.Bytes()))

		listener.Close()
	}()

	// Finally, just listen for file operations
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err)
			break
		}

		go rpc.ServeConn(conn)
	}
}
