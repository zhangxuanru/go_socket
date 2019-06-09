package main

import (
	"socket/app/protocol/server"
)

func main() {
	 server.Init()
}

//	var (
//		tcpAddr   *net.TCPAddr
//		tcpListen *net.TCPListener
//		tcpConn   *net.TCPConn
//		err       error
//	)
//	if tcpAddr, err = net.ResolveTCPAddr(common.GetNetWorkType(), common.GetServerAddress()); err != nil {
//		fmt.Println("ResolveTCPAddr error:", err)
//		return
//	}
//	if tcpListen, err = net.ListenTCP(common.GetNetWorkType(), tcpAddr); err != nil {
//		fmt.Println("ListenTCP error:", err)
//		return
//	}
//	for {
//		if tcpConn, err = tcpListen.AcceptTCP(); err != nil {
//			continue
//		}
//		go handleClient(tcpConn)
//	}
//}
//
//func handleClient(conn *net.TCPConn) {
//	defer conn.Close()
//	fmt.Println("conn:", conn.RemoteAddr().String())
//	var (
//		buf []byte
//	)
//	buf = make([]byte, 1024)
//	for {
//		conn.Read(buf)
//		fmt.Println(string(buf))
//	}
//}
