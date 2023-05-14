package client

import (
	"log"
	"math/rand"
	"net"
	"os"
	"time"

	"github.com/namelew/RPC/packages/messages"
)

const protocol = "tcp"

type Client struct {
	id      int
	adress  string
	usetime time.Duration
}

func New(id int, adress string, useDelay time.Duration) *Client {
	return &Client{
		id:      id,
		adress:  adress,
		usetime: useDelay,
	}
}

func (c *Client) Run() {
	for {
		if c.Lock() {
			log.Printf("Client %d enter the critical region\n", c.id)
			time.Sleep(c.usetime)
			c.Unlock()
		} else {
			waitTime := rand.Intn(10)
			log.Printf("Client %d can't enter the critical region and will sleep %d seconds\n", c.id, waitTime)
			time.Sleep(time.Second * time.Duration(waitTime))
		}
	}
}

func (c *Client) Lock() bool {
	conn, err := net.Dial(protocol, os.Getenv("CTRADRESS"))
	buffer := make([]byte, 256)

	if err != nil {
		log.Println("Unable to create connection with Coordenator. ", err.Error())
		return false
	}

	request := messages.Message{
		Action:   messages.REQUEST,
		Lockback: c.adress,
	}

	requestPayload, err := request.Pack()

	if err != nil {
		log.Println("Unable to create request message. ", err.Error())
		return false
	}

	conn.Write(requestPayload)

	time.After(time.Second * 5)

	n, err := conn.Read(buffer)

	if err != nil {
		log.Println("Unable to read data from coordinator. ", err.Error())
		return false
	}

	response := messages.Message{}

	response.Unpack(buffer[:n])

	return response.Action == messages.ALLOW
}

func (c *Client) Unlock() {
	conn, err := net.Dial(protocol, os.Getenv("CTRADRESS"))

	if err != nil {
		log.Println("Unable to create connection with Coordenator. ", err.Error())
		return
	}

	request := messages.Message{
		Action:   messages.FREE,
		Lockback: c.adress,
	}

	requestPayload, err := request.Pack()

	if err != nil {
		log.Println("Unable to create unlock message. ", err.Error())
		return
	}

	conn.Write(requestPayload)

	time.After(time.Second * 5)
}
