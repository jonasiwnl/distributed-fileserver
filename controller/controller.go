package controller

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"time"
)

const (
	PORT = ":2121"
)

type FileServer struct {
	Addr     string
	SizeUsed int64
	Capacity int64
}

type FileServerArgs struct {
	SizeUsed int64
	Capacity int64
}

type FileMetadata struct {
	Name       string
	Size       int64
	Created    time.Time
	Modified   time.Time
	ChunkAddrs []string
}

type Controller struct {
	FileServers []FileServer
}

func (c *Controller) listenForFileServers(quit chan bool) {
	addr := net.UDPAddr{
		Port: 9999,
		IP:   net.ParseIP("0.0.0.0"),
	}

	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	defer conn.Close()

	buffer := make([]byte, 1024)

	for {
		select {
		case <-quit:
			return
		default:
			n, remoteAddr, err := conn.ReadFromUDP(buffer)
			if err != nil {
				fmt.Printf("Failed to read UDP message: %v", err)
				continue
			}

			var fileServer FileServerArgs
			decoder := gob.NewDecoder(bytes.NewBuffer(buffer[:n]))
			err = decoder.Decode(&fileServer)
			if err != nil {
				fmt.Printf("Failed to decode data: %v", err)
				continue
			}

			c.FileServers = append(c.FileServers, FileServer{remoteAddr.String(), fileServer.SizeUsed, fileServer.Capacity})
			fmt.Println("Registering fileserver at address: ", remoteAddr.String())
		}
	}
}

func Start(quit chan bool) {
	controller := new(Controller)
	rpc.Register(controller)

	listener, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Printf("Couldn't listen on port %s: %s\n", PORT, err)
		return
	}
	defer listener.Close()

	fmt.Printf("Controller listening on port %s\n", PORT)

	go controller.listenForFileServers(quit)

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
