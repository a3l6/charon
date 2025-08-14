package main

import (
	"bufio"
	"fmt"
	"net"
)

func runClient(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	/* GENERAL MESSAGE SENDING
	- write the frame, check for errors
	- flush the writer

	GENERAL MESSAGE RECIEVING
	- read frame and pass reader, check for errors
	*/

	if err := writeFrame(writer, NewFilename, payload); err != nil {

	}
}
