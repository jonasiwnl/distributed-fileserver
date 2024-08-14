package server

import (
	"io/fs"
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

type DirArgs struct {
	Path string
	Mode os.FileMode
}

type ListArgs struct {
	Path string
}

type FileInfo struct {
	Name    string
	Size    int64
	Mode    os.FileMode
	ModTime time.Time
}

type DirEntry struct {
	Name  string
	IsDir bool
	Type  fs.FileMode
	Info  FileInfo
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

func (fs *FileServer) ListDirectory(args *ListArgs, reply *[]DirEntry) error {
	virtualPath := getVirtualPath(args.Path)
	files, err := os.ReadDir(virtualPath)
	if err != nil {
		return err
	}
	// Convert fs.DirEntry to gob encoding friendly DirEntry.
	for _, file := range files {
		info, _ := file.Info()
		*reply = append(*reply, DirEntry{
			Name:  file.Name(),
			IsDir: file.IsDir(),
			Type:  info.Mode(),
			Info: FileInfo{
				Name:    file.Name(),
				Size:    info.Size(),
				Mode:    info.Mode(),
				ModTime: info.ModTime(),
			},
		})
	}
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
