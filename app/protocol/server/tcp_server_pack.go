package server

import (
	"net"
	"errors"
)

type serverPack struct {
}

func (s *serverPack) Write(conn *net.TCPConn, data []byte) (int, error) {
	if conn == nil {
		return 0, errors.New("conn is nil")
	}
	return conn.Write(data)
}

func NewServerPack() *serverPack {
	 return &serverPack{}
}
