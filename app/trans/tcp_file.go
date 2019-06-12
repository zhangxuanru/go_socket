package trans

import (
	"errors"
	"fmt"
	"net"
	"socket/app/config"
)

type TcpFilePack struct {
	closeChan chan bool
}

func (s *TcpFilePack) Write(conn *net.TCPConn, data []byte) (int, error) {
	if conn == nil {
		return 0, errors.New("conn is nil")
	}
	data = append([]byte(config.SEND_FILE_HEADER_PACK), data...)
	return conn.Write(Pack(data))
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

func NewTcpFilePack() *TcpFilePack {
	return &TcpFilePack{
		closeChan: make(chan bool),
	}
}

func (s *TcpFilePack) ClosePack() {
	s.closeChan <- true
}
