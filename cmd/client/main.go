package main

import (
	"flag"
	"log"
	"math/rand"
	"time"

	"github.com/joho/godotenv"
	"github.com/namelew/RPC/internal/client"
)

const USETIME = time.Second * 2

// TODO:
//	- Critical region ins't be granted
//  - Don't now if is client ou coordinator

func main() {
	var (
		id = flag.Int("id", rand.Int(), "client id")
	)
	godotenv.Load()

	for {
		if client.Request() {
			log.Printf("Client %d enter the critical region\n", *id)
			time.Sleep(USETIME)
		} else {
			waitTime := rand.Int()
			log.Printf("Client %d can't enter the critical region and will sleep %d seconds\n", *id, waitTime)
			time.Sleep(time.Second * time.Duration(waitTime))
		}
	}
}
