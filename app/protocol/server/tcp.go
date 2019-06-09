package server

import (
	"net"
	"socket/app/common"
	"fmt"
	"socket/app/protocol/io"
	"socket/app/service"
)

type server struct {
	*service.ClientMsg
	servicePck *service.ServerPack
}

func Init()  {
    s := &server{
		ClientMsg:service.NewClientMsg(),
    	servicePck:service.NewServerPack(),
	}
    if err := s.Listen();err!=nil{
    	fmt.Println("error:",err)
    	return
	}
}

func (s *server) Listen() error {
	var(
		tcpAddr *net.TCPAddr
		tcpListen *net.TCPListener
		tcpConn *net.TCPConn
		err error
	)
    if tcpAddr,err = net.ResolveTCPAddr(common.GetNetWorkType(),common.GetServerAddress());err!=nil{
    	goto PrintErr
	}
	if tcpListen,err = net.ListenTCP(common.GetNetWorkType(),tcpAddr);err!=nil{
		goto PrintErr
	}
	fmt.Println("ListenTCP:",tcpAddr.String()," start")
	go io.NewStdInIo().OutStdInMsgByChan(s.InputMsgChan)
	defer tcpListen.Close()
	for{
		if tcpConn,err = tcpListen.AcceptTCP();err!=nil{
			fmt.Println("error:",err)
			continue
		}
       go s.handleClient(tcpConn)
	}
	return err

	PrintErr:
		fmt.Println("error:",err)
	return err
}

func (s *server) handleClient(conn *net.TCPConn) {
     defer func(conn *net.TCPConn) {
     	 if conn!=nil{
			 conn.Close()
		 }
	 }(conn)
     s.PrintClientConnMsg(conn)
     for{
		 select {
		 case s.InputMsg = <- s.InputMsgChan:
                s.servicePck.Write(conn,[]byte(s.InputMsg))
		}
	 }
}

func (s *server) PrintClientConnMsg(conn *net.TCPConn)  {
      fmt.Println(conn.RemoteAddr().String())
}

