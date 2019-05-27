package main

import (
	"net"
	"fmt"
	"bufio"
	"strings"
	"time"
)
//TCP服务器，带读取超时的判断
func main()  {
	startServer()
}

func startServer()  {
	var(
		listener net.Listener
		conn net.Conn
		err error
	)
	if listener,err = net.Listen("tcp",":2000");err!=nil{
        fmt.Println("TCP服务启动失败，失败原因:",err.Error())
        return
	}
	fmt.Println("TCP服务启动成功---")
	fmt.Println("服务器信息,Network:",listener.Addr().Network(),"  String:",listener.Addr().String())
	defer listener.Close()
	for{
	    if conn,err = listener.Accept();err!=nil{
			fmt.Println("TCP服务读取数据失败，失败原因:",err.Error())
			break
		}
		//conn.SetDeadline()
		//go handleConnection(conn)
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	var (
		timer       *time.Timer
		readTcpData chan bool
	)
	readTcpData = make(chan bool)
	timer = time.NewTimer(5 * time.Second)
	go handleConnection(conn, readTcpData)
	for {
		select {
		case <-timer.C:
			fmt.Println("read msg time out-----")
			goto RETRY
		case c := <-readTcpData:
			fmt.Println("readTcpData channel-----", c)
			timer.Reset(time.Duration(time.Second * 5))
		}
	}
 RETRY:
    fmt.Println("close handleClient")
	if conn !=nil{
		conn.Close()
	}
}

func handleConnection(conn net.Conn,readTcpData chan bool)  {
	var(
		readMsg string
		err error
	)
	fmt.Println("Server conn RemoteAddr:",conn.RemoteAddr().String(),"  LocalAddr:",conn.LocalAddr().String())
	for{
		if conn == nil{
			fmt.Println("无效的 socket 连接")
			break
		}
	    if readMsg,err =  bufio.NewReader(conn).ReadString('\n');err!=nil{
           fmt.Println("读取消息错误:",err.Error())
           break
		}
		readMsg = strings.TrimSpace(readMsg)
		if strings.EqualFold(readMsg,"Bye"){
			fmt.Println("close conn")
            break
		}
		if _,err := conn.Write([]byte(fmt.Sprintf("HI，字符串长度是:%d\r\n",len(readMsg))));err!=nil{
			fmt.Println("send error,",err.Error())
		}
		fmt.Println("send OK")
		readTcpData<-true
	}
	fmt.Println("conn close..........")
	if conn != nil{
	    conn.Close()
	}
}




