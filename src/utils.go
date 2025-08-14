package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io"
	"math"
)

func sumArrBytes(obj []byte) int {
	var res int

	for _, val := range obj {
		res += int(val)
	}

	return res
}

func writeFile(file File) error {
	fmt.Println("Writing file")
	return nil
}

// Write frame writes a frame according to the Charon spec. This frame can be sent to another instance of Charon and be decoded properly.
func writeFrame(w io.Writer, msgType Message, payload []byte) error {
	buf := &bytes.Buffer{}

	// write headers
	if _, err := buf.Write(Magic); err != nil {
		return err
	}
	if err := buf.WriteByte(Version); err != nil {
		return err
	}
	if err := buf.WriteByte(byte(msgType)); err != nil {
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

// Readframe reads a given frame (r) and outputs the result (in order): message type, version, payload, error
func readFrame(r io.Reader) (Message, byte, []byte, error) {
	header := make([]byte, HeaderLen)

	if _, err := io.ReadFull(r, header); err != nil {
		return 0, 0, nil, err
	}

	// validate
	if !bytes.Equal(header[0:4], Magic) {
		return 0, 0, nil, BadMagic
	}

	ver := header[4]
	msgType := Message(header[5])

	length := binary.BigEndian.Uint32(header[6:10])

	const twentyMegabytes uint32 = 10 << 20
	if length > twentyMegabytes {
		return 0, 0, nil, fmt.Errorf("payload of size %d too large", length)
	}

	payload := make([]byte, length)
	if _, err := io.ReadFull(r, payload); err != nil {
		return 0, 0, nil, err
	}
	var crc uint32
	if err := binary.Read(r, binary.BigEndian, &crc); err != nil {
		return 0, 0, nil, err
	}
	if crc != crc32.ChecksumIEEE(payload) {
		return 0, 0, nil, BadCRC
	}

	return msgType, ver, payload, nil
}
