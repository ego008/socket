package socket

import (
	"bytes"
	"io"
	"net"
	"sync"
	"time"
)

// Socket ...
type Socket struct {
	net.Conn
	isUse bool
	mu    sync.RWMutex
}

// Release ...
func (s *Socket) Release() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.isUse = false

}

// IsUse ...
func (s *Socket) IsUse() (b bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	b = s.isUse
	return
}

// Use ...
func (s *Socket) Use() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.isUse = true
}

// ReadAll ...
func (s *Socket) ReadAll(initialCap ...int) (datas []byte, err error) {
	initial := 4096
	if len(initialCap) > 0 && initialCap[0] > 0 {
		initial = initialCap[0]
	}
	request := make([]byte, initial)
	var jsonBuf bytes.Buffer
	var n int
	for {
		n, err = s.Read(request)
		if n > 0 {
			jsonBuf.Write(request[0:n])
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
	}
	return jsonBuf.Bytes(), nil
}

// Ping to detect whether the socket is closed.
func (s *Socket) Ping() (b bool, err error) {
	one := []byte{}
	s.SetReadDeadline(time.Now())
	if _, err := s.Read(one); err == io.EOF {
		s.Close()
		return false, err
	}
	var zero time.Time
	s.SetReadDeadline(zero)
	return true, err
}
