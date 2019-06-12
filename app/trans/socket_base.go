package trans

import (
	"errors"
	"fmt"
	"io"
	"net"
	"socket/app/common"
)

type SocketSendBase interface {
	Write(conn *net.TCPConn, data []byte) (int, error)
	Read(conn *net.TCPConn, readMsgChan chan []byte) error
	Receive(receiveMsg []byte)
	Close(conn *net.TCPConn) error
}

type ServerTcpBase interface {
	Write(conn *net.TCPConn, data []byte) (int, error)
}

type SocketBase struct {
	TcpAddr *net.TCPAddr
	TcpConn *net.TCPConn
}

type ClientMsg struct {
	InputMsg          string
	InputMsgChan      chan string
	ReadMsg           []byte
	ReadMsgChan       chan []byte
	IsCloseChan       chan bool
	IsCloseServerChan chan bool
}

func NewClientMsg() *ClientMsg {
	return &ClientMsg{
		InputMsgChan:      make(chan string),
		ReadMsgChan:       make(chan []byte),
		IsCloseChan:       make(chan bool),
		IsCloseServerChan: make(chan bool),
	}
}

func NewSocketBase(addr *net.TCPAddr, conn *net.TCPConn) SocketBase {
	return SocketBase{
		TcpAddr: addr,
		TcpConn: conn,
	}
}

func TcpRead(conn *net.TCPConn, readChan chan []byte) (err error) {
	var (
		tmpBuf []byte
		buf    []byte
	)
	if conn == nil {
		return errors.New("conn is nil")
	}
	tmpBuf = make([]byte, 0)
	unpackChan := make(chan []byte)
	for {
		if conn == nil {
			goto CLOSE
		}
		buf = make([]byte, 1024)
		_, err = conn.Read(buf)
		if err != nil {
			goto CLOSE
		}
		go func() {
			for {
				select {
				case packBuf := <-unpackChan:
					if common.IsBufBye(packBuf) {
						err = io.EOF
						goto CLOSE
					}
					readChan <- packBuf
				}
			}
		CLOSE:
		}()
		//解包
		tmpBuf = UnPack(append(tmpBuf, buf...), unpackChan)
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
