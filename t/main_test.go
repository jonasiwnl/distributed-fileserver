package distributed_fileserver_test

import (
	"net/rpc"
	"testing"

	"github.com/jonasiwnl/distributed-fileserver/v2/server"
)

const (
	CONTROLLERPORT = ":2120"
	FILESERVERPORT = ":2125"
)

// Share global clients for all tests.
var FileServerClient *rpc.Client
var ControllerClient *rpc.Client

func TestMain(m *testing.M) {
	fscQuit := make(chan bool, 1)
	csQuit := make(chan bool, 1)
	sixtyFourMB := int64(64 * 1024 * 1024)
	go server.StartControllerServer(CONTROLLERPORT, csQuit)
	go server.StartFileServer(FILESERVERPORT, sixtyFourMB, fscQuit)

	var fscErr, csErr error
	FileServerClient, fscErr = rpc.Dial("tcp", "localhost"+FILESERVERPORT)
	ControllerClient, csErr = rpc.Dial("tcp", "localhost"+CONTROLLERPORT)

	if fscErr == nil && csErr == nil {
		// Run tests
		m.Run()
	}

	// Cleanup
	if fscErr == nil {
		FileServerClient.Close()
	}
	if csErr == nil {
		ControllerClient.Close()
	}
	fscQuit <- true
	csQuit <- true
}
