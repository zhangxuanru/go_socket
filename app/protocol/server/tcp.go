package server

import (
	"net"
	"fmt"
	"socket/app/protocol/io"
	"socket/app/socket"
	"socket/app/common"
	"runtime"
)

type server struct {
	*socket.ClientMsg
    *socket.SocketIo
	tcpAddr *net.TCPAddr
	tcpListen *net.TCPListener
}

func Init()  {
    s := &server{
		ClientMsg: socket.NewClientMsg(),
	}
    if err := s.Listen();err!=nil{
    	return
	}
}

func (s *server) Listen() error {
	var(
		tcpConn *net.TCPConn
		err error
	)
    if s.tcpAddr,err = net.ResolveTCPAddr(common.GetNetWorkType(),common.GetServerAddress());err!=nil{
    	goto PrintErr
	}
	if s.tcpListen,err = net.ListenTCP(common.GetNetWorkType(),s.tcpAddr);err!=nil{
		goto PrintErr
	}
	fmt.Println("ListenTCP:",s.tcpAddr.String()," start")
	go io.NewStdInIo().OutStdInMsgByChan(s.InputMsgChan)
	defer s.tcpListen.Close()
	for{
		if tcpConn,err = s.tcpListen.AcceptTCP();err!=nil{
			fmt.Println("error:",err)
			continue
		}
       go s.handleClient(tcpConn)
	}
	return nil
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
	 sockConnBase := socket.NewSocketBase(s.tcpAddr,conn)
	 s.SocketIo = socket.NewSocketIo(sockConnBase)
	 go s.SocketIo.SocketPack.Read(conn,s.ReadMsgChan)
     for{
		 select {
		 case s.InputMsg = <- s.InputMsgChan:
		 	fmt.Println("--server input--:")
               s.SocketIo.SocketPack.Write(conn,[]byte(s.InputMsg))
		       s.SocketIo.CheckBye(s.InputMsg,s.IsCloseChan)
		 case s.ReadMsg = <-s.ReadMsgChan:
		 	fmt.Println("---server read :--")
		 	   s.SocketIo.SocketPack.Receive(s.ReadMsg)
			   s.SocketIo.CheckBye(string(s.ReadMsg),s.IsCloseChan)
		 case <-s.IsCloseChan:
		 	goto END
		}
	}
	END:
		fmt.Println("conn end:",conn.RemoteAddr().String())
	    s.SocketIo.SocketPack.Close(conn)
}

func (s *server) PrintClientConnMsg(conn *net.TCPConn)  {
      fmt.Println("client:",conn.RemoteAddr().String(),"maxgoroutine:",runtime.NumGoroutine())
}

