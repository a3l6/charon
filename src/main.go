package main

import (
	"errors"
	"fmt"
	"log"
	"net"
)

// Protocol constants
var Magic = []byte{'c', 'h', 'a', 'r', 'o', 'n'}

const Version byte = 1

type Message byte

const (
	Text Message = iota
	Data
	NewFilePayload
	NewFileSize
	NewFilename
)

const HeaderLen = 4 + 1 + 1 + 4 // magic(4)+version(1)+type(1)+length(4)

var BadMagic = errors.New("bad magic")
var BadCRC = errors.New("bad crc32")

type File struct {
	setname  bool
	name     string
	contents []byte
}

func main() {
	addr, err := net.ResolveTCPAddr("tcp", ":8000")
	if err != nil {
		log.Fatal(err)
	}

	ln, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	defer ln.Close()

	fmt.Println("Listening on port 8000")
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConnection(conn)
	}
}
