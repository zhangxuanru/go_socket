package trans

import (
	"fmt"
	"net"
)

type TcpFilePack struct {
}

func (s *TcpFilePack) Write(conn *net.TCPConn, data []byte) (int, error) {
	return 0, nil
}

func (s *TcpFilePack) Read(conn *net.TCPConn, readMsgChan chan []byte) (err error) {
	return nil
}

func (s *TcpFilePack) Close(conn *net.TCPConn) error {
	return nil
}

func (s *TcpFilePack) Receive(receiveMsg []byte) {
	fmt.Println("receiveMsg:", receiveMsg)
}
