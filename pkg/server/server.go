package server

import (
	"bytes"
	"log"
	"net"
	"strings"
	"sync"
)

type HandlerFunc func(conn net.Conn)

type Server struct {
	addr     string
	mu       sync.RWMutex
	handlers map[string]HandlerFunc
}

func NewServer(addr string) *Server {
	return &Server{addr: addr, handlers: make(map[string]HandlerFunc)}
}

func (s *Server) Register(path string, handler HandlerFunc) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.handlers[path] = handler
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Print(err)
		return err
	}
	defer func() {
		if cerr := listener.Close(); cerr != nil {
			if err == nil {
				err = cerr
				return
			}
			log.Print(cerr)
		}
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			return err

		}

		err = s.handle(conn)
		if err != nil {
			log.Print(err)
			continue
		}
	}

}

func (s Server) handle(conn net.Conn) (err error) {
	defer func() {
		if cerr := conn.Close(); cerr != nil {
			if err == nil {
				err = cerr
				return
			}
			log.Print(err)
		}
	}()

	buf := make([]byte, 4096)

	n, err := conn.Read(buf)
	data := buf[:n]
	requestedLineDelim := []byte{'\r', '\n'}
	requestedLineEnd := bytes.Index(data, requestedLineDelim)
	if requestedLineEnd == -1 {
		return err
	}

	requestLine := string(data[:requestedLineEnd])
	parts := strings.Split(requestLine, " ")
	if len(parts) != 3 {
		return err
	}

	_, path, version := parts[0], parts[1], parts[2]

	if version != "HTTP/1.1" {
		return
	}

	if handler, exist := s.handlers[path]; exist {
		handler(conn)
	}

	return nil
}
