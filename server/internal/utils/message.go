package utils

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

func ReadMsg(conn net.Conn) ([]byte, error) {
	var length uint64
	err := binary.Read(conn, binary.BigEndian, &length)
	if err != nil {
		return nil, fmt.Errorf("failed to read msg: %w", err)
	}

	buf := make([]byte, length)
	_, err = io.ReadFull(conn, buf)
	if err != nil {
		return nil, fmt.Errorf("failed to read buf: %w", err)
	}

	return buf, nil
}

func WriteMsg(conn net.Conn, data []byte) error {
	length := uint64(len(data))
	if err := binary.Write(conn, binary.BigEndian, length); err != nil {
		return fmt.Errorf("failed to write length: %w", err)
	}

	n, err := conn.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write data: %w", err)
	}
	if n != len(data) {
		return fmt.Errorf("incomplete write: wrote %d of %d bytes", n, len(data))
	}

	return nil
}
