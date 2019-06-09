package service

import (
	"strings"
)

const SEND_STRING = 1
const SEND_FILE = 2

type SocketIo struct {
	SocketBase
	SocketPack SocketSendBase
	SendType   int //1:字符串，2：文件
}

func (s *SocketIo) initSockPack() {
	if s.SocketPack == nil {
		s.SocketPack = &stringPack{}
		s.SendType = SEND_STRING
	}
}

func (s *SocketIo) WriteData(msg string) {
	if strings.HasPrefix(msg, "file:") {
		s.SocketPack = &filePack{}
		s.SendType = SEND_FILE
	}
	s.SocketPack.Write(s.TcpConn, []byte(msg))
}

func (s *SocketIo) CheckBye (msg string,IsCloseChan chan bool)  {
    if strings.EqualFold(msg,"bye"){
		IsCloseChan<-true
	}
}


func NewSocketIo(sockBase SocketBase) *SocketIo {
	socketIo:= &SocketIo{
		SocketBase: sockBase,
	}
	socketIo.initSockPack()
    return socketIo
}



