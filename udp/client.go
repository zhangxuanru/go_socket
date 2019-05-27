package main

import (
	"net"
	"fmt"
	"bufio"
	"os"
)

func main()  {
	connectUdpServer()
}

func connectUdpServer()  {
	var(
		udpAddr *net.UDPAddr
		conn *net.UDPConn
		err error
	)
     if udpAddr,err = net.ResolveUDPAddr("udp","127.0.0.1:8100");err!=nil{
     	fmt.Println("ResolveUDPAddr error:",err.Error())
     	return
	 }
	 if conn,err = net.DialUDP("udp",nil,udpAddr);err!=nil{
	 	fmt.Println("DialUDP error:",err.Error())
	 	return
	 }
	 sendMsg(conn)
}

func sendMsg(conn *net.UDPConn)  {
	var(
		msg string
		err error
	)
	defer conn.Close()
	for{
		fmt.Println("input msg:")
		if msg,err = bufio.NewReader(os.Stdin).ReadString('\n');err!=nil{
			fmt.Println("ReadString error:",err.Error())
			return
		}
		conn.Write([]byte(msg))
		if msg,err = bufio.NewReader(conn).ReadString('\n');err!=nil{
			fmt.Println("ReadString error:",err.Error())
			break
		}
		fmt.Println("read msg:",msg)
	}
}
