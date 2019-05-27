package main

import (
	"net"
	"fmt"
	"bytes"
)

func main()  {
	startUdpServer()
}

func startUdpServer()  {
	var(
		udpAddr *net.UDPAddr
		udpConn *net.UDPConn
		err error
	)
	if udpAddr,err = net.ResolveUDPAddr("udp","127.0.0.1:8100");err!=nil{
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
		buf []byte
		msg []byte
		udpAddr *net.UDPAddr
		err error
	)
	buf = make([]byte,1024)
	msg = make([]byte,1024)
	fmt.Println("RemoteAddr:",conn.RemoteAddr()," LocalAddr:",conn.LocalAddr())
	if _,udpAddr,err = conn.ReadFromUDP(buf);err!=nil{
		fmt.Println("ReadFromUDP error:",err.Error())
	}
	fmt.Println("RemoteAddr IP :",udpAddr.IP,"RemoteAddr port:",udpAddr.Port)
	msg = []byte(fmt.Sprintf("hi:%s",getStringByBuF2(buf)))
	if _,err = conn.WriteToUDP(msg,udpAddr);err!=nil{
		fmt.Println("WriteToUDP error:",err.Error())
		return
	}
	fmt.Println("input msg:",string(buf))
	fmt.Println("reverse msg:",string(msg))
}


func getStringByBuf(buf []byte) (string) {
	var str string
	for i,v := range buf{
		if v == 0{
			str = string(buf[:i])
			break
		}
	}
	return str
}

func getStringByBuF2(buf []byte) string {
	index := bytes.IndexByte(buf,0)
	return string(buf[:index])
}
