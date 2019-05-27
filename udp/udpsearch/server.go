package main

import (
	"net"
	"fmt"
	msg2 "socket/udp/udpsearch/msg"
	"bytes"
)

var SN =  "SN-ZXR1"
func main()  {
	startServer()
}

func startServer()  {
	fmt.Println("udp server start:")
	var(
		udpAddr *net.UDPAddr
		udpConn *net.UDPConn
		err error
	)
     if udpAddr,err = net.ResolveUDPAddr("udp",":20000");err!=nil{
     	fmt.Println("ResolveUDPAddr error:",err.Error())
     	return
	 }
	 if udpConn,err = net.ListenUDP("udp",udpAddr);err!=nil{
		 fmt.Println("ListenUDP error:",err.Error())
		 return
	 }
	 defer udpConn.Close()
	 for{
          handleClient(udpConn)
	 }
}

func handleClient(conn *net.UDPConn)  {
	var(
		msgStr string
		port int
		respMsg string
		err error
		buf []byte
		udpAddr *net.UDPAddr
	)
	buf = make([]byte,1024)
	if _,udpAddr,err = conn.ReadFromUDP(buf);err!=nil{
		fmt.Println("ReadFromUDP error:",err.Error())
		return
	}
	  msgStr = getStringByBuF2(buf)
	 if port,err = msg2.ParsePort(msgStr);err!=nil{
	 	fmt.Println("ParsePort error:",err.Error())
	 	return
	 }
	 fmt.Println("client port",port,"client:",udpAddr.String())
	 respMsg = msg2.BuildWithSn(SN)

	 //发送到指定端口
	 clientUdp := &net.UDPAddr{
	 	IP: []byte(udpAddr.IP),
	 	Port:port,
	 }

	 if _,err = conn.WriteToUDP([]byte(respMsg),clientUdp);err!=nil{
	 	fmt.Println("WriteToUDP error:",err.Error())
	 	return
	 }

	fmt.Println("client msg:",msgStr)
	fmt.Println("server msg:",respMsg)
}


func getStringByBuF2(buf []byte) string {
	index := bytes.IndexByte(buf,0)
	return string(buf[:index])
}
