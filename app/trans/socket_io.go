package trans

import (
	"socket/app/common"
	"strings"
)

type SocketIo struct {
	SocketBase
	SocketPack    SocketSendBase
	TransportType string
}

func (s *SocketIo) initSockPack() {
	s.TransportType = common.GetNetWorkType()
	if s.SocketPack == nil {
		s.SocketPack = &TcpStringPack{}
	}
}

func (s *SocketIo) WriteData(msg string) {
	if strings.HasPrefix(msg, "file:") {
		s.SocketPack = &TcpFilePack{}
	}
	s.SocketPack.Write(s.TcpConn, []byte(msg))
}

func NewSocketIo(sockBase SocketBase) *SocketIo {
	socketIo := &SocketIo{
		SocketBase: sockBase,
	}
	socketIo.initSockPack()
	return socketIo
}
