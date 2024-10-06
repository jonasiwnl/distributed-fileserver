package client

import "net/rpc"

type Client struct {
	controllerClient *rpc.Client
}

func (c *Client) AddFile(args struct{}, reply struct{}) error {
	return nil
}

func NewClient() *Client {
	controllerClient, err := rpc.Dial("tcp", "localhost:2120")
	if err != nil {
		panic(err)
	}

	client := Client{controllerClient}
	rpc.Register(client)
	return &client
}
