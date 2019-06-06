package service

import (
	"strings"
)

const SEND_STRING = 1
const SEND_FILE = 2

type SocketIo struct {
	SocketBase
	socketPack SocketSendBase
	SendType   int //1:字符串，2：文件
}

func (s *SocketIo) initSockPack() {
	if s.socketPack == nil {
		s.socketPack = &stringPack{}
		s.SendType = SEND_STRING
	}
}

func (s *SocketIo) WriteData(msg string) {
	s.initSockPack()
	if strings.HasPrefix(msg, "file:") {
		s.socketPack = &filePack{}
		s.SendType = SEND_FILE
	}
	s.socketPack.Write(s.TcpConn, []byte(msg))
}

func (s *SocketIo) Read() {
	s.initSockPack()
	s.socketPack.Read(s.TcpConn)
}

func (s *SocketIo) Close() {
	s.initSockPack()
	s.socketPack.Close(s.TcpConn)
}

func NewSocketIo(sockBase SocketBase) *SocketIo {
	return &SocketIo{
		SocketBase: sockBase,
	}
}
