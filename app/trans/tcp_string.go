package trans

import (
	"errors"
	"fmt"
	"io"
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
	var (
		tmpBuf []byte
		buf    []byte
	)
	if conn == nil {
		return errors.New("conn is nil")
	}
	tmpBuf = make([]byte, 0)
	for {
		if conn == nil {
			goto CLOSE
		}
		select {
		case <-s.closeChan:
			fmt.Println("---closeChan--")
			err = errors.New("auto close")
			goto CLOSE
		default:
			buf = make([]byte, 1024)
			if _, err = conn.Read(buf); err != nil {
				goto CLOSE
			}
			//解包
			tmpBuf = UnPack(append(tmpBuf, buf...), readChan)
		}
	}
	return err
CLOSE:
	tmpBuf = make([]byte, 0)
	buf = make([]byte, 0)
	fmt.Println("read close......")
	if err != nil && err == io.EOF {
		readChan <- []byte("bye")
	}
	return err
}

func (s *TcpStringPack) Close(conn *net.TCPConn) error {
	fmt.Println("conn is close")
	return conn.Close()
}

func (s *TcpStringPack) Receive(receiveMsg []byte) {
	if common.IsBufBye(receiveMsg) {
		go s.ClosePack()
	}
	fmt.Printf("receiveMsg:%s\n", string(receiveMsg))
}

func NewTcpStringPack() *TcpStringPack {
	return &TcpStringPack{
		closeChan: make(chan bool),
	}
}

func (s *TcpStringPack) ClosePack() {
	s.closeChan <- true
}
