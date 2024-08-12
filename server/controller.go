package server

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
	CONTROLLERPORT = ":2120"
)

type FileServerMessageType int

const (
	REGISTER FileServerMessageType = iota
	HEARTBEAT
	DISCONNECT
)

type FileServerData struct {
	SizeUsed int64
	Capacity int64
}

type FileServerMessage struct {
	Type FileServerMessageType
	Data FileServerData
}

type FileMetadata struct {
	Name       string
	Size       int64
	Created    time.Time
	Modified   time.Time
	ChunkAddrs []string
}

type Controller struct {
	FileServers map[string]FileServerData
}

func (c *Controller) GetFileServers(args struct{}, reply *map[string]FileServerData) error {
	*reply = c.FileServers
	return nil
}

// TODO
func (c *Controller) FindFile(args struct{}, reply *FileMetadata) error {
	return nil
}

// TODO
func (c *Controller) FindDir(args struct{}, reply *[]byte) error {
	return nil
}

func (c *Controller) listenForFileServers(quit chan bool) {
	addr := net.UDPAddr{
		Port: 9998,
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

			var fileServerMessage FileServerMessage
			decoder := gob.NewDecoder(bytes.NewBuffer(buffer[:n]))
			err = decoder.Decode(&fileServerMessage)
			if err != nil {
				fmt.Printf("Failed to decode data: %v", err)
				continue
			}

			switch fileServerMessage.Type {
			case REGISTER:
				c.FileServers[remoteAddr.String()] = fileServerMessage.Data
				fmt.Println("Registering fileserver at address: ", remoteAddr.String())
			case HEARTBEAT:
				// TODO: a heartbeat would be nice.
			case DISCONNECT:
				delete(c.FileServers, remoteAddr.String())
			}
		}
	}
}

func StartControllerServer(quit chan bool) {
	controller := new(Controller)
	rpc.Register(controller)

	listener, err := net.Listen("tcp", CONTROLLERPORT)
	if err != nil {
		fmt.Printf("Couldn't listen on port %s: %s\n", CONTROLLERPORT, err)
		return
	}
	defer listener.Close()

	fmt.Printf("Controller listening on port %s\n", CONTROLLERPORT)

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
