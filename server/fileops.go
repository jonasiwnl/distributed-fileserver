package server

import (
	"os"
	"path/filepath"
	"sync"
	"time"
)

var locks map[string]*sync.RWMutex = make(map[string]*sync.RWMutex)

func getVirtualPath(path string) string {
	return filepath.Join(DIRECTORY, path)
}

type FileArgs struct {
	Path string
	Data []byte
	Mode os.FileMode
}

type FileInfo struct {
	Name    string
	Size    int64
	Mode    os.FileMode
	ModTime time.Time
}

func (fs *FileServer) ReadFile(args *FileArgs, reply *[]byte) error {
	virtualPath := getVirtualPath(args.Path)
	locks[virtualPath].RLock()
	defer locks[virtualPath].RUnlock()
	data, err := os.ReadFile(virtualPath)
	*reply = data
	return err
}

func (fs *FileServer) WriteFile(args *FileArgs, reply *bool) error {
	virtualPath := getVirtualPath(args.Path)
	if _, ok := locks[virtualPath]; ok {
		locks[virtualPath].Lock()
		defer locks[virtualPath].Unlock()
	} else {
		locks[virtualPath] = &sync.RWMutex{}
	}
	err := os.WriteFile(virtualPath, args.Data, args.Mode)
	*reply = err == nil
	return err
}

func (fs *FileServer) RemoveFile(args *FileArgs, reply *bool) error {
	virtualPath := getVirtualPath(args.Path)
	locks[virtualPath].Lock()
	err := os.Remove(virtualPath)
	*reply = err == nil
	locks[virtualPath].Unlock()
	delete(locks, args.Path)
	return err
}
