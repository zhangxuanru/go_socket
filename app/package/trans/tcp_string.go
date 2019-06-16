package trans

import (
	"errors"
	"fmt"
	"net"
	"socket/app/common"
	"socket/app/config"
)

type TcpStringPack struct {
	closeChan chan bool
}

func (s *TcpStringPack) Write(conn *net.TCPConn, data []byte) (int, error) {
	if conn == nil {
		return 0, errors.New("conn is nil")
	}
	data = append([]byte(config.SEND_STR_HEADER_PACK), data...)
	return conn.Write(Pack(data))
}

func (s *TcpStringPack) Read(conn *net.TCPConn, readChan chan []byte) (err error) {
	NewRead().Read(conn, readChan)
	return nil
}

func (s *TcpStringPack) Close(conn *net.TCPConn) error {
	fmt.Println("conn is close")
	return conn.Close()
}

func (s *TcpStringPack) Receive(receiveMsg []byte) {
	if common.IsBufBye(receiveMsg) {
		NewRead().Close()
	}
	fmt.Printf("receiveMsg:%s\n", string(receiveMsg))
}

func NewTcpStringPack() *TcpStringPack {
	return &TcpStringPack{
		closeChan: make(chan bool),
	}
}
