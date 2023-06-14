package main

import (
	"log"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", ":3000")

	if err != nil {
		log.Fatal(err)
	}

	_, err = conn.Write([]byte("SET foo bar123 5000000"))

	if err != nil {
		log.Fatal(err)
	}
}
