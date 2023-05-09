package client

import (
	"log"
	"net"
	"os"
	"time"

	"github.com/namelew/RPC/packages/messages"
)

const protocol = "tcp"

func Request() bool {
	c, err := net.Dial(protocol, os.Getenv("CTRADRESS"))
	buffer := make([]byte, 256)

	if err != nil {
		log.Println("Unable to create connection with Coordenator. ", err.Error())
		return false
	}

	request := messages.Message{
		Action: messages.REQUEST,
	}

	requestPayload, err := request.Pack()

	if err != nil {
		log.Println("Unable to create request message. ", err.Error())
		return false
	}

	c.Write(requestPayload)

	time.After(time.Second * 5)

	n, err := c.Read(buffer)

	if err != nil {
		log.Println("Unable to read data from coordinator. ", err.Error())
		return false
	}

	response := messages.Message{}

	response.Unpack(buffer[:n])

	return response.Action == messages.ALLOW
}
