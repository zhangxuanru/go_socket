package main

import (
	"net"
	"fmt"
	"time"
)

func main() {
	// 这里设置发送者的IP地址，自己查看一下自己的IP自行设定
	laddr := net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 20000,
	}
	// 这里设置接收者的IP地址为广播地址
	raddr := net.UDPAddr{
		IP:   net.IPv4(255, 255, 255, 255),
		Port: 30201,
	}
	//conn, err := net.ListenUDP("udp",&laddr)
	conn, err := net.DialUDP("udp", &laddr, &raddr)
	if err != nil {
		println(err.Error())
		return
	}
	conn.Write([]byte("7,7,7,7,7,7,7这是暗号,请回电商品(PORT):20000"))


	buf := make([]byte,1024)
	conn.ReadFromUDP(buf)
	fmt.Println(string(buf))
	for{
		time.Sleep(1 * time.Second)
	}
	fmt.Println("-----------")
//	conn.Close()
}

