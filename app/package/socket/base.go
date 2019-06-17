package socket

import (
	"net"
	"os"
)

type SocketSendBaseer interface {
	Write(conn *net.TCPConn, data []byte) (int, error)
	Read(conn *net.TCPConn, readMsgChan chan []byte) error
	Receive(receiveMsg []byte, receiveChan chan []byte, source string)
	SocketSendCloseer
}

type SocketSendCloseer interface {
	Close(conn *net.TCPConn) error
}

type ServerTcpBaseer interface {
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
	ReceiveChan       chan []byte
	Signal            chan os.Signal
}

func NewClientMsg() *ClientMsg {
	return &ClientMsg{
		InputMsgChan:      make(chan string),
		ReadMsgChan:       make(chan []byte),
		IsCloseChan:       make(chan bool),
		IsCloseServerChan: make(chan bool),
		ReceiveChan:       make(chan []byte),
		Signal:            make(chan os.Signal),
	}
}

func NewSocketBase(addr *net.TCPAddr, conn *net.TCPConn) SocketBase {
	return SocketBase{
		TcpAddr: addr,
		TcpConn: conn,
	}
}
