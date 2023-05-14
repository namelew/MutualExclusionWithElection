package main

import (
	"flag"
	"math/rand"
	"time"

	"github.com/joho/godotenv"
	"github.com/namelew/RPC/internal/client"
)

const USETIME = time.Second * 2

// TODO:
//	- Implementing control queue to requests

func main() {
	var (
		id     = flag.Int("id", rand.Int(), "client id")
		adress = flag.String("adress", "localhost:30002", "client lockback adress")
	)
	godotenv.Load()

	c := client.New(*id, *adress, USETIME)
	c.Run()
}
