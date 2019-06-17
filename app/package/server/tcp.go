package server

import (
	"fmt"
	"net"
	"runtime"
	"socket/app/common"
	"socket/app/config"
	"socket/app/package/io"
	"socket/app/package/socket"
)

type server struct {
	*socket.ClientMsg
	*socket.SocketIo
	tcpAddr   *net.TCPAddr
	tcpListen *net.TCPListener
}

func Init() {
	s := &server{
		ClientMsg: socket.NewClientMsg(),
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
	sockConnBase := socket.NewSocketBase(s.tcpAddr, conn)
	s.SocketIo = socket.NewSocketIo(sockConnBase, socket.SEND_STRING)
	go s.SocketIo.ReadData(conn, s.ReadMsgChan)

	for {
		select {
		case s.InputMsg = <-s.InputMsgChan:
			s.SocketIo.SocketPack.Write(conn, []byte(s.InputMsg))
			common.CheckBye([]byte(s.InputMsg), s.IsCloseServerChan)
		case s.ReadMsg = <-s.ReadMsgChan:
			if s.SocketIo.IsStrByCurrPack(s.ReadMsg) == false {
				s.SocketIo.ResetSocketPack()
			}
			s.ReadMsg = common.RemoveStrSendHeader(s.ReadMsg)
			s.SocketIo.SocketPack.Receive(s.ReadMsg, s.ReceiveChan, config.SERVERIDENT)
			common.CheckBye(s.ReadMsg, s.IsCloseServerChan)
		case receiveMsg := <-s.ReceiveChan:
			fmt.Println(string(receiveMsg))
			s.SocketIo.SocketPack.Write(conn, receiveMsg)
		case <-s.IsCloseServerChan:
			fmt.Println("close server")
			goto END
		}
	}
	return
END:
	fmt.Println("conn end:", conn.RemoteAddr().String())
	s.SocketIo.SocketPack.Close(conn)
	fmt.Println("end handleClient.......", "--maxgoroutine:", runtime.NumGoroutine())
	return
}

func (s *server) PrintClientConnMsg(conn *net.TCPConn) {
	fmt.Println("client:", conn.RemoteAddr().String(), "maxgoroutine:", runtime.NumGoroutine())
}
