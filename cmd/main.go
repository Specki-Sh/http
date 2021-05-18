package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/Specki-Sh/http/pkg/server"
)

func main() {
	host := "localhost"
	port := "9999"

	if err := execute(host, port); err != nil {
		os.Exit(1)

	}
	fmt.Println("server closed")
}

func execute(host string, port string) (err error) {
	srv := server.NewServer(net.JoinHostPort(host, port))
	srv.Register("/", func(req *server.Request) {
		body := "Welcome to our web-site"
		id := req.QueryParams["id"]
		log.Print(id)

		_, err = req.Conn.Write([]byte(
			"HTTP/1.1 200 OK\r\n" +
				"Content-Lenght: " + strconv.Itoa(len(body)) + "\r\n" +
				"Content-Type: text/html\r\n" +
				"Connection: close\r\n" +
				"\r\n" +
				body,
		))

		if err != nil {
			log.Print(err)
		}
	})

	srv.Register("/payment/{id}", func(req *server.Request) {
		id := req.PathParams["id"]
		log.Print(id)
	})

	srv.Register("/about", func(req *server.Request) {
		body := "About Golang Academy"

		_, err = req.Conn.Write([]byte(
			"HTTP/1.1 200 OK\r\n" +
				"Content-Lenght: " + strconv.Itoa(len(body)) + "\r\n" +
				"Content-Type: text/html\r\n" +
				"Connection: close\r\n" +
				"\r\n" +
				body,
		))

		if err != nil {
			log.Print(err)
		}
	})
	return srv.Start()
}
