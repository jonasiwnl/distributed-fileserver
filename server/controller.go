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

type FileServerMessageType int

const (
	REGISTER FileServerMessageType = iota
	HEARTBEAT
	DISCONNECT
)

type FileServerEntry struct {
	Addr string
	Data FileServerData
}

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
	AddrToIdx   map[string]int
	FileServers []FileServerEntry
}

func NewController() *Controller {
	return &Controller{
		AddrToIdx:   make(map[string]int),
		FileServers: make([]FileServerEntry, 0),
	}
}

func (c *Controller) GetFileServers(args struct{}, reply *[]FileServerEntry) error {
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
				fmt.Println("Registering fileserver at address: ", remoteAddr.String())
				c.FileServers = append(c.FileServers, FileServerEntry{Addr: remoteAddr.String(), Data: fileServerMessage.Data})
				c.AddrToIdx[remoteAddr.String()] = len(c.FileServers) - 1
			case HEARTBEAT:
				// TODO: a heartbeat would be nice.
			case DISCONNECT:
				fmt.Println("Disconnecting fileserver at address: ", remoteAddr.String())
				// Find index of fileserver to remove.
				idx := c.AddrToIdx[remoteAddr.String()]
				// Swap the last element to the index.
				c.FileServers[idx] = c.FileServers[len(c.FileServers)-1]
				// Update the index of the swapped element.
				c.AddrToIdx[c.FileServers[idx].Addr] = idx
				// Remove the last element.
				c.FileServers = c.FileServers[:len(c.FileServers)-1]
				delete(c.AddrToIdx, remoteAddr.String())
			}
		}
	}
}

func StartControllerServer(port string, quit chan bool) {
	controller := NewController()
	rpc.Register(controller)

	listener, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Printf("Couldn't listen on port %s: %s\n", port, err)
		return
	}
	defer listener.Close()

	fmt.Printf("Controller listening on port %s\n", port)

	go controller.listenForFileServers(quit)

	go func() {
		<-quit
		listener.Close()
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err)
			break
		}

		go rpc.ServeConn(conn)
	}
}
