package socket_test

import (
	"errors"
	"log"
	"net"
	"testing"
	"time"

	"github.com/mushroomsir/socket"
	"github.com/stretchr/testify/assert"
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
func TestSocket(t *testing.T) {
	t.Run("Socket use default options that should be", func(t *testing.T) {
		assert := assert.New(t)
		pool, err := socket.NewPool(1, 20, socket.NewFactory(address))
		if !assert.Nil(err) {
			return
		}
		assert.Equal(20, pool.Cap())
		assert.Equal(1, pool.Len())
		socket, err := pool.Get()
		if !assert.Nil(err) {
			return
		}
		b, err := socket.Ping()
		if !assert.Nil(err) {
			return
		}
		assert.Equal(true, b)
	})
	t.Run("Socket use error options that should be", func(t *testing.T) {
		assert := assert.New(t)
		pool, err := socket.NewPool(1, 0, socket.NewFactory(address))
		assert.Nil(pool)
		assert.Equal(socket.ErrCapacitySettings, err)

		pool, err = socket.NewPool(1, 20, func() (conn net.Conn, err error) {
			return nil, errors.New("failed")
		})
		if assert.NotNil(err) {
			assert.Contains(err.Error(), "factory is not able to fill the pool")
		}
		pool, err = socket.NewPool(1, 20, socket.NewFactory(address))
		pool.Close()
		conn, err := pool.Get()
		assert.Nil(conn)
		assert.Equal(socket.ErrClosed, err)

		pool, err = socket.NewPool(1, 5, socket.NewFactory(address))
		for index := 0; index < pool.Cap(); index++ {
			pool.Get()
		}
		conn, err = pool.Get()
		assert.Nil(conn)
		assert.Equal(socket.ErrMaxedOut, err)

		i := 0
		pool, err = socket.NewPool(1, 20, func() (conn net.Conn, err error) {
			i++
			t.Log(i)
			if i == 1 {
				return socket.NewFactory(address)()
			}
			return nil, errors.New("failed")
		})
		conn, err = pool.Get()
		conn, err = pool.Get()
		// assert.Nil(conn)
		// if assert.NotNil(err) {
		// 	assert.Contains(err.Error(), "factory is not able to fill the pool")
		// }

	})
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
