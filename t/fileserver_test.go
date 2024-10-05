// As this is a side project, these tests aren't super extensive.
// Just checking that there aren't horrible errors.

package distributed_fileserver_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jonasiwnl/distributed-fileserver/v2/server"
)

func TestAddRemoveFile(t *testing.T) {
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
