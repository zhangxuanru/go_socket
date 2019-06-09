package service

import (
	"net"
	"fmt"
)

type filePack struct {
}

func (s *filePack) Write(conn *net.TCPConn, data []byte) (int, error) {
	return 0, nil
}

func (s *filePack) Read(conn *net.TCPConn,readMsgChan chan  []byte) (err error) {
	return  nil
}

func (s *filePack) Close(conn *net.TCPConn) error {
	return nil
}

func (s *filePack) Receive(receiveMsg []byte)  {
	fmt.Println("receiveMsg:",receiveMsg)
}
