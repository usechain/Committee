package sssa

import (
	"fmt"
	"gitlab.com/usechain/go-usechain/commitee/gsocket"
	"time"
)

type demoClient struct{}

func (client *demoClient) OnConnect(c *gsocket.Connection) {
	fmt.Printf("CONNECTED: %s\n", c.RemoteAddr())
}

func (client *demoClient) OnDisconnect(c *gsocket.Connection) {
	fmt.Printf("DISCONNECTED: %s\n", c.RemoteAddr())
}

func (client *demoClient) OnRecv(c *gsocket.Connection, data []byte) {
	//fmt.Printf("DATA RECVED: %s %d - %v\n", c.RemoteAddr(), len(data), data)
}

func (client *demoClient) OnError(c *gsocket.Connection, err error) {
	fmt.Printf("ERROR: %s - %s\n", c.RemoteAddr(), err.Error())
}

func SendMsg(port uint16, msg []byte) (string, error) {
	demoClient := &demoClient{}

	client := gsocket.CreateTCPClient(demoClient.OnConnect, demoClient.OnDisconnect, demoClient.OnRecv, demoClient.OnError)

	err := client.Connect("127.0.0.1", port)
	if err != nil {
		fmt.Printf("Connect Server Error: %s\n", err.Error())
		return "", err
	}
	fmt.Printf("Connect Server %s Success\n", client.RemoteAddr())
	client.Send(msg)
	time.Sleep(1000000)
	client.Close()
	return  "okay", nil
}


