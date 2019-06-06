package service

import "net"

type SocketSendBase interface {
	Write(conn *net.TCPConn, data []byte) (int, error)
	Read(conn *net.TCPConn) ([]byte, error)
	Close(conn *net.TCPConn) error
}

type SocketBase struct {
	TcpAddr *net.TCPAddr
	TcpConn *net.TCPConn
}

type ClientMsg struct {
	InputMsg     string
	InputMsgChan chan string
}

func NewClientMsg() ClientMsg {
	return ClientMsg{
		InputMsgChan: make(chan string),
	}
}
