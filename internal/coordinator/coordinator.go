package coordinator

import (
	"bufio"
	"log"
	"net"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"github.com/namelew/RPC/packages/messages"
	"github.com/namelew/RPC/packages/queue"
)

const protocol = "tcp"

type Coordinator struct {
	request        *queue.Queue[string]
	mutex          *sync.Mutex
	criticalRegion bool
}

func Build() *Coordinator {
	return &Coordinator{
		request:        &queue.Queue[string]{},
		mutex:          &sync.Mutex{},
		criticalRegion: false,
	}
}

func (cd *Coordinator) Handler() {
	godotenv.Load(".env")

	l, err := net.Listen(protocol, os.Getenv("CTRADRESS"))

	if err != nil {
		log.Panic("Unable to create lisntener. ", err.Error())
	}

	for {
		request, err := l.Accept()

		if err != nil {
			log.Println("Unable to serve connection. ", err.Error())
			continue
		}

		go func(c net.Conn) {
			var in, out messages.Message
			b := make([]byte, 1024)

			n, err := bufio.NewReader(c).Read(b)

			defer c.Close()

			if err != nil {
				log.Println("Unable to read data from "+c.RemoteAddr().String()+". ", err.Error())
				return
			}

			if err := in.Unpack(b[:n]); err != nil {
				log.Println("Unable to unpack data from "+c.RemoteAddr().String()+". ", err.Error())
				return
			}

			switch in.Action {
			case messages.REQUEST:
				cd.mutex.Lock()
				defer cd.mutex.Unlock()

				if !cd.criticalRegion {
					cd.criticalRegion = true
					log.Println(c.RemoteAddr().String(), "allowed to access critical region")
					out.Action = messages.ALLOW
				} else {
					log.Println(c.RemoteAddr().String(), "not allowed to access critical region")
					out.Action = messages.REFUSE
				}
			case messages.FREE:
				cd.mutex.Lock()
				defer cd.mutex.Unlock()
				cd.criticalRegion = false
				out.Action = messages.ACKFREE
				log.Println("Critical region free to use")
			}

			Send(c, &out)
		}(request)
	}
}

func Send(c net.Conn, m *messages.Message) {
	payload, err := m.Pack()

	if err != nil {
		log.Println("Unable to pack message. ", err.Error())
		return
	}

	_, err = c.Write(payload)

	if err != nil {
		log.Println("Unable to send message to "+c.RemoteAddr().String()+". ", err.Error())
		return
	}
}
