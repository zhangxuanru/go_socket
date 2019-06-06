package service

import (
	"errors"
	"net"
)

type stringPack struct {
}

func (s *stringPack) Write(conn *net.TCPConn, data []byte) (int, error) {
	if conn == nil {
		return 0, errors.New("conn is nil")
	}
	return conn.Write(data)
}

func (s *stringPack) Read(conn *net.TCPConn) (data []byte, err error) {
	if conn == nil {
		return nil, errors.New("conn is nil")
	}
	buf := make([]byte, 1024)
	if _, err = conn.Read(buf); err != nil {
		return nil, err
	}
	return buf, nil
}

func (s *stringPack) Close(conn *net.TCPConn) error {
	return nil
}
