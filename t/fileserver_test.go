// As this is a side project, these tests aren't super extensive.
// Just checking that there aren't horrible errors.

package distributed_fileserver_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jonasiwnl/distributed-fileserver/v2/server"
)

func TestDir(t *testing.T) {
	testDir := "testdir"
	directoryPath := filepath.Join(server.DIRECTORY, testDir)

	args := &server.DirArgs{Path: testDir, Mode: 0755}
	var reply bool
	err := FileServerClient.Call("FileServer.MakeDirectory", args, &reply)
	if err != nil {
		t.Fatal("making directory:", err)
	}
	if _, err := os.Stat(directoryPath); os.IsNotExist(err) || !reply {
		t.Fatal("directory not created")
	}

	err = FileServerClient.Call("FileServer.RemoveDirectory", args, &reply)
	if err != nil {
		t.Fatal("removing directory:", err)
	}
	if _, err := os.Stat(directoryPath); err == nil || !reply {
		t.Fatal("directory not removed")
	}

	err = FileServerClient.Call("FileServer.RemoveDirectory", args, &reply)
	if err != nil {
		t.Fatal("error removing non-existent directory:", err)
	}
}

func TestFile(t *testing.T) {
	testFile := "testfile"
	filePath := filepath.Join(server.DIRECTORY, testFile)

	args := &server.FileArgs{Path: testFile, Data: []byte("test"), Mode: 0644}
	var reply bool
	err := FileServerClient.Call("FileServer.WriteFile", args, &reply)
	if err != nil {
		t.Fatal("writing file:", err)
	}
	if _, err := os.Stat(filePath); os.IsNotExist(err) || !reply {
		t.Fatal("file not created")
	}

	err = FileServerClient.Call("FileServer.RemoveFile", args, &reply)
	if err != nil {
		t.Fatal("removing file:", err)
	}
	if _, err := os.Stat(filePath); err == nil || !reply {
		t.Fatal("file not removed")
	}
}

func TestListDir(t *testing.T) {
	testDir1 := "testdir1"
	testDir2 := "testdir2"
	testFile1 := "testfile1"
	testFile2 := "testfile2"
	testFile3 := "testfile3"

	var mkDirReply bool

	testDirArgs1 := &server.DirArgs{Path: testDir1, Mode: 0755}
	FileServerClient.Call("FileServer.MakeDirectory", testDirArgs1, &mkDirReply)
	testDirArgs2 := &server.DirArgs{Path: testDir2, Mode: 0755}
	FileServerClient.Call("FileServer.MakeDirectory", testDirArgs2, &mkDirReply)
	testFileArgs1 := &server.FileArgs{Path: filepath.Join(testDir1, testFile1), Data: []byte("test1"), Mode: 0644}
	FileServerClient.Call("FileServer.WriteFile", testFileArgs1, &mkDirReply)
	testFileArgs2 := &server.FileArgs{Path: filepath.Join(testDir1, testFile2), Data: []byte("test2"), Mode: 0644}
	FileServerClient.Call("FileServer.WriteFile", testFileArgs2, &mkDirReply)
	testFileArgs3 := &server.FileArgs{Path: testFile3, Data: []byte("test3"), Mode: 0644}
	FileServerClient.Call("FileServer.WriteFile", testFileArgs3, &mkDirReply)

	var listDirReply []server.DirEntry
	listDirArgs := &server.ListArgs{Path: ""}
	err := FileServerClient.Call("FileServer.ListDirectory", listDirArgs, &listDirReply)
	if err != nil {
		t.Fatal("listing directory:", err)
	}
	if len(listDirReply) != 3 {
		t.Fatal("incorrect number of entries")
	}

	testDir1Found := false
	testDir2Found := false
	testFile3Found := false
	for _, entry := range listDirReply {
		if entry.Name == testDir1 {
			testDir1Found = true
		}
		if entry.Name == testDir2 {
			testDir2Found = true
		}
		if entry.Name == testFile3 {
			testFile3Found = true
		}
	}
	if !testDir1Found || !testDir2Found || !testFile3Found {
		t.Fatal("missing entry.")
	}

	listDirArgs = &server.ListArgs{Path: testDir1}
	err = FileServerClient.Call("FileServer.ListDirectory", listDirArgs, &listDirReply)
	if err != nil {
		t.Fatal("listing directory:", err)
	}
	if len(listDirReply) != 2 {
		t.Fatal("incorrect number of entries")
	}

	testFile1Found := false
	testFile2Found := false

	for _, entry := range listDirReply {
		if entry.Name == testFile1 {
			testFile1Found = true
		}
		if entry.Name == testFile2 {
			testFile2Found = true
		}
	}
	if !testFile1Found || !testFile2Found {
		t.Fatal("missing entry.")
	}

	// Clean up files
	FileServerClient.Call("FileServer.RemoveDirectory", testDirArgs1, &mkDirReply)
	FileServerClient.Call("FileServer.RemoveDirectory", testDirArgs2, &mkDirReply)
	FileServerClient.Call("FileServer.RemoveFile", testFileArgs3, &mkDirReply)
}
