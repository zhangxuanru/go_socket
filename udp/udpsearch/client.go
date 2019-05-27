package main

import (
	"net"
	"fmt"
	"socket/udp/udpsearch/msg"
	"bytes"
)

func main()  {
    connUdpServer()
	startUdpListenServer()
}

func startUdpListenServer()  {
	fmt.Println("start listen server:")
	var(
		udpAddr *net.UDPAddr
		udpConn *net.UDPConn
		err error
	)
	 if udpAddr,err = net.ResolveUDPAddr("udp",":5000");err!=nil{
	 	fmt.Println("ResolveUDPAddr client error:",err.Error())
	 	return
	 }
	 if udpConn,err = net.ListenUDP("udp",udpAddr);err!=nil{
	 	fmt.Println("ListenUDP client error:",err.Error())
	 	return
	 }
	 defer udpConn.Close()
	 for{
         handleServer(udpConn)
	  }
}

func handleServer(conn *net.UDPConn)  {
	var(
		buf []byte
		err error
		udpAddr *net.UDPAddr
	)
	buf = make([]byte,1024)
	if _,udpAddr,err = conn.ReadFromUDP(buf);err!=nil{
		fmt.Println("ReadFromUDP error:",err.Error())
		return
	}
	fmt.Println("server ip:",udpAddr.IP," server port:",udpAddr.Port)
	fmt.Println("read msg:",getStringByBuF3(buf))
}


func connUdpServer()  {
	var(
		udpAddr *net.UDPAddr
		udpConn *net.UDPConn
		err error
	)
   if udpAddr,err =  net.ResolveUDPAddr("udp","255.255.255.255:20000");err!=nil{
		fmt.Println("ResolveUDPAddr error:",err.Error())
		return
   }
   if udpConn,err = net.DialUDP("udp",nil,udpAddr);err!=nil{
		fmt.Println("DialUDP error:",err.Error())
		return
   }
   defer udpConn.Close()
   sendMsg(udpConn)
}

func sendMsg(conn *net.UDPConn)  {
	  var(
	  	respMsg  string
	  	err error
	  )
	 respMsg = msg.BuildWithPort(5000)
	 if _,err =conn.Write([]byte(respMsg));err!=nil{
	 	fmt.Println("Write error:",err.Error())
	 	return
	 }
	 fmt.Println("send ok")
}


func getStringByBuF3(buf []byte) string {
	index := bytes.IndexByte(buf,0)
	return string(buf[:index])
}
