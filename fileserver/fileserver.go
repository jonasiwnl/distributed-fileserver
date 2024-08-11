package fileserver

import (
	"fmt"
	"io/fs"
	"net"
	"net/rpc"
	"os"
	"path/filepath"
)

const (
	PORT      = ":2122"
	DIRECTORY = "virtual/"
)

func getVirtualPath(path string) string {
	return filepath.Join(DIRECTORY, path)
}

// ******** FILE STUFF ******** //

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
	virtualPath := getVirtualPath(args.Path)
	data, err := os.ReadFile(virtualPath)
	*reply = data
	return err
}

func (fs *FileServer) WriteFile(args *FileArgs, reply *bool) error {
	virtualPath := getVirtualPath(args.Path)
	err := os.WriteFile(virtualPath, args.Data, args.Mode)
	*reply = err == nil
	return err
}

func (fs *FileServer) RemoveFile(args *FileArgs, reply *bool) error {
	virtualPath := getVirtualPath(args.Path)
	err := os.Remove(virtualPath)
	*reply = err == nil
	return err
}

func (fs *FileServer) ListDirectory(args *ListArgs, reply *[]fs.DirEntry) error {
	virtualPath := getVirtualPath(args.Path)
	files, err := os.ReadDir(virtualPath)
	if err != nil {
		return err
	}
	*reply = files
	return nil
}

func (fs *FileServer) MakeDirectory(args *DirArgs, reply *bool) error {
	virtualPath := getVirtualPath(args.Path)
	err := os.MkdirAll(virtualPath, args.Mode)
	*reply = err == nil
	return err
}

func (fs *FileServer) RemoveDirectory(args *DirArgs, reply *bool) error {
	virtualPath := getVirtualPath(args.Path)
	err := os.RemoveAll(virtualPath)
	*reply = err == nil
	return err
}

// ******** FILE STUFF ******** //

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
