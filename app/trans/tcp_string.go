package trans

import (
	"errors"
	"fmt"
	"io"
	"net"
	"socket/app/common"
)

type TcpStringPack struct {
}

func (s *TcpStringPack) Write(conn *net.TCPConn, data []byte) (int, error) {
	if conn == nil {
		return 0, errors.New("conn is nil")
	}
	return conn.Write(data)
}

func (s *TcpStringPack) Read(conn *net.TCPConn, readChan chan []byte) (err error) {
	if conn == nil {
		return errors.New("conn is nil")
	}
	for {
		if conn == nil {
			goto CLOSE
		}
		buf := make([]byte, 1024)
		if _, err = conn.Read(buf); err != nil {
			goto CLOSE
		}
		if common.IsBufBye(buf) {
			err = io.EOF
			goto CLOSE
		}
		readChan <- buf
	}
	return err
CLOSE:
	fmt.Println("read close error:", err)
	if err != nil && err == io.EOF {
		readChan <- []byte("bye")
	}
	return err
}

func (s *TcpStringPack) Close(conn *net.TCPConn) error {
	return conn.Close()
}

func (s *TcpStringPack) Receive(receiveMsg []byte) {
	fmt.Printf("receiveMsg:%s\n", string(receiveMsg))
}
