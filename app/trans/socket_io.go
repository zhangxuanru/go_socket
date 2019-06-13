package trans

import (
	"socket/app/common"
	"socket/app/config"
	"strings"
)

const SEND_STRING = 1
const SEND_FILE = 2

type SocketIo struct {
	SocketBase
	SocketPack    SocketSendBase
	TransportType string
	SendType      int
}

func (s *SocketIo) initSockPack(sendType int) {
	s.TransportType = common.GetNetWorkType()
	if sendType == SEND_STRING {
		s.SetStringPack()
	}
	if sendType == SEND_FILE {
		s.SetFilePack()
	}
	if s.SocketPack == nil {
		s.SetStringPack()
	}
}

func (s *SocketIo) WriteData(msg string) {
	msg = strings.TrimSpace(msg)
	if len(msg) == 0 {
		return
	}
	if strings.HasPrefix(msg, "file:") {
		s.SetFilePack()
	}
	s.SocketPack.Write(s.TcpConn, []byte(msg))
}

func (s *SocketIo) SetFilePack() {
	s.SocketPack = NewTcpFilePack()
	s.SendType = SEND_FILE
}

func (s *SocketIo) SetStringPack() {
	s.SocketPack = NewTcpStringPack()
	s.SendType = SEND_STRING
}

//检查当前的消息是否是在当前正确的pack中接收的
func (s *SocketIo) IsStrByCurrPack(str []byte) bool {
	strPackLen := len(config.SEND_STR_HEADER_PACK)
	if len(str) < strPackLen {
		return true
	}
	leftMsg := string(str[0:strPackLen])
	if leftMsg != config.SEND_STR_HEADER_PACK && leftMsg != config.SEND_FILE_HEADER_PACK {
		return true
	}
	if leftMsg == config.SEND_STR_HEADER_PACK && s.SendType == SEND_STRING {
		return true
	}
	if leftMsg == config.SEND_FILE_HEADER_PACK && s.SendType == SEND_FILE {
		return true
	}
	return false
}

func (s *SocketIo) ResetSocketPack() {
	if s.SendType == SEND_STRING {
		s.SetFilePack()
	} else {
		s.SetStringPack()
	}
}

func NewSocketIo(sockBase SocketBase, sendType int) *SocketIo {
	socketIo := &SocketIo{
		SocketBase: sockBase,
	}
	socketIo.initSockPack(sendType)
	return socketIo
}
