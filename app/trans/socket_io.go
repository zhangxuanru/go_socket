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
		s.SocketPack = NewTcpStringPack()
	}
}

func (s *SocketIo) WriteData(msg string) {
	if strings.HasPrefix(msg, "file:") {
		s.SocketPack = NewTcpFilePack()
	}
	s.SocketPack.Write(s.TcpConn, []byte(msg))
}

func (s *SocketIo) SetFilePack() {
	s.SocketPack = NewTcpFilePack()
}

/*
这个方法需要优化，可以加个字段根据 字段类型来判断
func (s *SocketIo) IsStrByWherePack(str []byte) bool {
	fmt.Println(reflect.TypeOf(s.SocketPack).String())
	filePackLen := len(config.SEND_FILE_HEADER_PACK)
	if string(str[0:filePackLen]) == config.SEND_FILE_HEADER_PACK && reflect.TypeOf(s.SocketPack).String() == "*trans.TcpFilePack" {
		return true
	}
	return false
}
*/

func NewSocketIo(sockBase SocketBase) *SocketIo {
	socketIo := &SocketIo{
		SocketBase: sockBase,
	}
	socketIo.initSockPack()
	return socketIo
}
