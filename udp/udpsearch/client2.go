package main

import (
	"net"
	"fmt"
	"socket/udp/udpsearch/msg"
)

func main()  {
	connUdp2()
}

func connUdp2()  {
	var(
		udpAddr *net.UDPAddr
		localUdp *net.UDPAddr
		udpConn *net.UDPConn
		err error
	)
	if udpAddr,err =  net.ResolveUDPAddr("udp","255.255.255.255:20000");err!=nil{
		fmt.Println("ResolveUDPAddr error:",err.Error())
		return
	}
	if localUdp,err = net.ResolveUDPAddr("udp",":5000");err!=nil{
		fmt.Println("ResolveUDPAddr local error:",err.Error())
		return
	}

	if udpConn,err = net.ListenUDP("udp",localUdp);err!=nil{
		fmt.Println("local ListenUDP error:",err.Error())
		return
	}
	defer udpConn.Close()
	send(udpConn,udpAddr)
}

func send(conn *net.UDPConn,srcUdp *net.UDPAddr)  {
	var(
		udpAddr *net.UDPAddr
		respMsg  string
		buf []byte
		err error
	)
	buf = make([]byte,1024)
	respMsg = msg.BuildWithPort(5000)
	if _,err = conn.WriteToUDP([]byte(respMsg),srcUdp);err!=nil{
		fmt.Println("WriteToUDP error:",err.Error())
		return
	}
	if _,udpAddr,err = conn.ReadFromUDP(buf);err!=nil{
		fmt.Println("ReadFromUDP error:",err.Error())
		return
	}
	fmt.Println("IP:",udpAddr.IP,"PORT:",udpAddr.Port)
	fmt.Println(string(buf))
}
