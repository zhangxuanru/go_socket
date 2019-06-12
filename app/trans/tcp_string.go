package trans

import (
	"errors"
	"fmt"
	"net"
)

type TcpStringPack struct {
}

func (s *TcpStringPack) Write(conn *net.TCPConn, data []byte) (int, error) {
	if conn == nil {
		return 0, errors.New("conn is nil")
	}
	return conn.Write(Pack(data))
}

func (s *TcpStringPack) Read(conn *net.TCPConn, readChan chan []byte) (err error) {
	err = TcpRead(conn, readChan)
	return err
}

func (s *TcpStringPack) Close(conn *net.TCPConn) error {
	fmt.Println("conn is close")
	return conn.Close()
}

func (s *TcpStringPack) Receive(receiveMsg []byte) {
	fmt.Printf("receiveMsg:%s\n", string(receiveMsg))
}
