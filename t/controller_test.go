package distributed_fileserver_test

import (
	"net"
	"testing"
	"time"

	"github.com/jonasiwnl/distributed-fileserver/v2/server"
)

func TestGetFileServers(t *testing.T) {
	// This timeout controls flakiness of the test.
	// Longer timeouts = less flakiness, but slower tests.
	// Since these tests are run locally, low values should be good.
	dialerWithTimeout := net.Dialer{Timeout: 1 * time.Millisecond}
	testFileServerPort := ":2126"

	getFileServersResponse := []server.FileServerEntry{}
	err := ControllerClient.Call("Controller.GetFileServers", struct{}{}, &getFileServersResponse)
	if err != nil {
		t.Fatal("getting file servers: ", err)
	}

	if len(getFileServersResponse) != 1 {
		t.Fatal("expected 1 file server baseline, got", len(getFileServersResponse))
	}

	quit := make(chan bool, 1)
	go server.StartFileServer(testFileServerPort, 0, quit)

	// We connect to our test file server. Doing this ensures that this testing thread waits
	// until the file server has connected to the controller before testing the number of file servers.
	testFileServerClient, err := dialerWithTimeout.Dial("tcp", "localhost"+testFileServerPort)
	if err != nil {
		t.Fatal("dialing file server: ", err)
	}

	err = ControllerClient.Call("Controller.GetFileServers", struct{}{}, &getFileServersResponse)

	// Quit here, so even if these tests fail, the fileserver still quits.
	testFileServerClient.Close()
	quit <- true

	if err != nil {
		t.Fatal("getting file servers: ", err)
	}

	if len(getFileServersResponse) != 2 {
		t.Fatal("expected 2 file servers after connect, got", len(getFileServersResponse))
	}

	// We try to connect to the file server again. This again ensures that this testing thread waits
	// until the file server has disconnected from the controller before testing the number of file servers.
	_, _ = dialerWithTimeout.Dial("tcp", "localhost"+testFileServerPort)

	// Test that the quit <- true actually worked.
	err = ControllerClient.Call("Controller.GetFileServers", struct{}{}, &getFileServersResponse)
	if err != nil {
		t.Fatal("getting file servers: ", err)
	}

	if len(getFileServersResponse) != 1 {
		t.Fatal("expected 1 file server after disconnect, got", len(getFileServersResponse))
	}
}
