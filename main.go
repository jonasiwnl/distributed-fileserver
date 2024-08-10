package main

import (
	"fmt"
	"net"
	"net/rpc"
	"os"
	"path/filepath"
)

const (
	PORT      = ":2122"
	DIRECTORY = "fileserver/"
)

func getCanonicalPath(path string) string {
	return filepath.Join(DIRECTORY, path)
}

type FileServer struct{}

type FileArgs struct {
	Path string
	Data []byte
	Mode os.FileMode // ?
}

type DirArgs struct {
	Path string
	Mode os.FileMode // ?
}

type ListArgs struct {
	Path string
}

func (fs *FileServer) ReadFile(args *FileArgs, reply *[]byte) error {
	canonical := getCanonicalPath(args.Path)
	data, err := os.ReadFile(canonical)
	*reply = data
	return err
}

func (fs *FileServer) WriteFile(args *FileArgs, reply *bool) error {
	canonical := getCanonicalPath(args.Path)
	err := os.WriteFile(canonical, args.Data, args.Mode)
	*reply = err == nil
	return err
}

func (fs *FileServer) RemoveFile(args *FileArgs, reply *bool) error {
	canonical := getCanonicalPath(args.Path)
	err := os.Remove(canonical)
	*reply = err == nil
	return err
}

func (fs *FileServer) MakeDirectory(args *DirArgs, reply *bool) error {
	canonical := getCanonicalPath(args.Path)
	err := os.MkdirAll(canonical, args.Mode)
	*reply = err == nil
	return err
}

func (fs *FileServer) RemoveDirectory(args *DirArgs, reply *bool) error {
	canonical := getCanonicalPath(args.Path)
	err := os.RemoveAll(canonical)
	*reply = err == nil
	return err
}

func main() {
	fileServer := new(FileServer)
	rpc.Register(fileServer)

	listener, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Printf("Couldn't listen on port %s: %s\n", PORT, err)
		return
	}

	fmt.Printf("Listening on port %s\n", PORT)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err)
			continue
		}

		go rpc.ServeConn(conn)
	}
}
