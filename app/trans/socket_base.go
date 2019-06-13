package trans

import (
	"net"
)

type SocketSendBase interface {
	Write(conn *net.TCPConn, data []byte) (int, error)
	Read(conn *net.TCPConn, readMsgChan chan []byte) error
	Receive(receiveMsg []byte)
	SocketSendClose
}

type SocketSendClose interface {
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
