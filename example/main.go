package main

import (
	"errors"
	"log"
	"net"
	"time"

	"github.com/mushroomsir/socket"
)

var (
	InitialCap = 5
	MaximumCap = 30
	network    = "tcp"
	address    = "127.0.0.1:7779"
	factory    = func() (net.Conn, error) { return net.Dial(network, address) }
)

func init() {
	// used for factory function
	go simpleTCPServer()
	time.Sleep(time.Millisecond * 300) // wait until tcp server has been settled
}
func main() {
	i := 0
	pool, _ := socket.NewPool(1, 20, func() (conn net.Conn, err error) {
		i++

		if i == 1 {
			return socket.NewFactory(address)()
		}
		return nil, errors.New("failed")
	})
	pool.Get()
	pool.Get()
}
func simpleTCPServer() {
	l, err := net.Listen(network, address)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		buffer := make([]byte, 256)
		conn.Read(buffer)
	}
}
