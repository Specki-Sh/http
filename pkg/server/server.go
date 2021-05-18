package main

import (
	"log"
	"net"
	"os"
	"strconv"

	"github.com/ripol92/http/pkg/server"
)

const CRLF = "\r\n"

func main() {
	host := "0.0.0.0"
	port := "9999"

	if err := execute(host, port); err != nil {
		os.Exit(1)
	}
}

func execute(host string, port string) (err error) {
	srv := server.NewServer(net.JoinHostPort(host, port))
	srv.Register("/", func(conn net.Conn) {
		body := "Welcome to our web-site"

		_, err = conn.Write([]byte(
			"HTTP/1.1 200 OK" + CRLF +
				"Content-Length: " + strconv.Itoa(len(body)) + CRLF +
				"Content-Type: text/html" + CRLF +
				"Connection: close" + CRLF +
				"\r\n" +
				body,
		))
		if err != nil {
			log.Print(err)
		}
	})
	srv.Register("/about", func(conn net.Conn) {
		body := "About Golang Academy"

		_, err = conn.Write([]byte(
			"HTTP/1.1 200 OK" + CRLF +
				"Content-Length: " + strconv.Itoa(len(body)) + CRLF +
				"Content-Type: text/html" + CRLF +
				"Connection: close" + CRLF +
				"\r\n" +
				body,
		))
		if err != nil {
			log.Print(err)
		}
	})
	return srv.Start()
}
