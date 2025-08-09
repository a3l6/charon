package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"log"
	"net"
)

// Protocol constants
var Magic = []byte{'c', 'h', 'a', 'r', 'o', 'n'}

const Version byte = 1

const (
	Text = iota
	Data
	File
)

const HeaderLen = 4 + 1 + 1 + 4 // magic(4)+version(1)+type(1)+length(4)

var BadMagic = errors.New("bad magic")
var BadCRC = errors.New("bad crc32")

func writeFrame(w io.Writer, msgType byte, payload []byte) error {
	buf := &bytes.Buffer{}

	// write headers
	if _, err := buf.Write(Magic); err != nil {
		return err
	}
	if err := buf.WriteByte(Version); err != nil {
		return err
	}
	if err := buf.WriteByte(msgType); err != nil {
		return err
	}

	// write length
	if err := binary.Write(buf, binary.BigEndian, uint32(len(payload))); err != nil {
		return err
	}

	// pay load
	// load: I can finally eat!
	if _, err := buf.Write(payload); err != nil {
		return err
	}

	// checksums
	// sum: Im good
	crc := crc32.ChecksumIEEE(payload)
	if err := binary.Write(buf, binary.BigEndian, crc); err != nil {
		return err
	}

	_, err := w.Write(buf.Bytes())
	return err
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Println(err)
			return
		}

		fmt.Println("Received: " + string(buf[:n]))
	}
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
