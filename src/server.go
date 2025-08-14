package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	// writer := bufio.NewWriter(conn)

	var file File

	for {
		msgType, version, payload, err := readFrame(reader)
		if err != nil {
			if err == io.EOF {
				writeFile(file)
				fmt.Println("client closed")
				return
			}
			fmt.Println("readframe error, closing client: ", err)
			return
		}

		switch version {
		case 1:
			switch msgType {
			case NewFilename:
				file.name = string(payload)
				file.setname = true
			case NewFilePayload:
				file.contents = append(file.contents, payload...)
			}

		}
	}
}
