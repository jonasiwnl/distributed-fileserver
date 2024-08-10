package main

import (
	"net/rpc"
	"os"
	"testing"
)

func TestMkRmDir(t *testing.T) {
	client, err := rpc.Dial("tcp", "localhost:2122")
	if err != nil {
		t.Fatal("dialing: ", err)
	}

	args := &DirArgs{Path: "testdir", Mode: 0755}
	var reply bool
	err = client.Call("FileServer.MakeDirectory", args, &reply)
	if err != nil {
		t.Fatal("making directory: ", err)
	}
	if _, err := os.Stat("fileserver/testdir"); os.IsNotExist(err) || !reply {
		t.Fatal("directory not created")
	}

	err = client.Call("FileServer.RemoveDirectory", args, &reply)
	if err != nil {
		t.Fatal("removing directory: ", err)
	}
	if _, err := os.Stat("fileserver/testdir"); err == nil || !reply {
		t.Fatal("directory not removed")
	}
}
