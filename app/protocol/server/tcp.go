package server

import (
	"fmt"
	"net"
	"runtime"
	"socket/app/common"
	"socket/app/protocol/io"
	"socket/app/trans"
)

type server struct {
	*trans.ClientMsg
	*trans.SocketIo
	tcpAddr   *net.TCPAddr
	tcpListen *net.TCPListener
}

func Init() {
	s := &server{
		ClientMsg: trans.NewClientMsg(),
	}
	if err := s.Listen(); err != nil {
		return
	}
}

func (s *server) Listen() error {
	var (
		tcpConn *net.TCPConn
		err     error
	)
	if s.tcpAddr, err = net.ResolveTCPAddr(common.GetNetWorkType(), common.GetServerAddress()); err != nil {
		goto PrintErr
	}
	if s.tcpListen, err = net.ListenTCP(common.GetNetWorkType(), s.tcpAddr); err != nil {
		goto PrintErr
	}
	fmt.Println("ListenTCP:", s.tcpAddr.String(), " start")
	go io.NewStdInIo(false).OutStdInMsgByChan(s.InputMsgChan)
	defer s.tcpListen.Close()
	for {
		if tcpConn, err = s.tcpListen.AcceptTCP(); err != nil {
			fmt.Println("error:", err)
			continue
		}
		go s.handleClient(tcpConn)
	}
	return nil
PrintErr:
	fmt.Println("error:", err)
	return err
}

func (s *server) handleClient(conn *net.TCPConn) {
	defer func(conn *net.TCPConn) {
		if conn != nil {
			conn.Close()
		}
	}(conn)
	s.PrintClientConnMsg(conn)
	sockConnBase := trans.NewSocketBase(s.tcpAddr, conn)
	s.SocketIo = trans.NewSocketIo(sockConnBase)
	go s.SocketIo.SocketPack.Read(conn, s.ReadMsgChan)

	for {
		select {
		case s.InputMsg = <-s.InputMsgChan:
			s.SocketIo.SocketPack.Write(conn, []byte(s.InputMsg))
			common.CheckBye([]byte(s.InputMsg), s.IsCloseServerChan)
		case s.ReadMsg = <-s.ReadMsgChan:
			// 这里去掉前缀 明天加
			//msg = bytes.Replace(msg, []byte(config.SEND_STR_HEADER_PACK), []byte(""), -1)
			//msg = bytes.Replace(msg, []byte(config.SEND_FILE_HEADER_PACK), []byte(""), -1)

			s.SocketIo.SocketPack.Receive(s.ReadMsg)
			common.CheckBye(s.ReadMsg, s.IsCloseServerChan)
		case <-s.IsCloseServerChan:
			fmt.Println("close server")
			goto END
		}
	}
END:
	fmt.Println("conn end:", conn.RemoteAddr().String())
	s.SocketIo.SocketPack.Close(conn)
	fmt.Println("end handleClient.......")
}

func (s *server) PrintClientConnMsg(conn *net.TCPConn) {
	fmt.Println("client:", conn.RemoteAddr().String(), "maxgoroutine:", runtime.NumGoroutine())
}
