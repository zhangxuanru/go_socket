package service

import "net"

type filePack struct {
}

func (s *filePack) Write(conn *net.TCPConn, data []byte) (int, error) {
	return 0, nil
}

func (s *filePack) Read(conn *net.TCPConn) (b []byte, err error) {
	return b, nil
}

func (s *filePack) Close(conn *net.TCPConn) error {
	return nil
}
