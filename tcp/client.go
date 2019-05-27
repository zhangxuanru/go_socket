package main

import (
	"fmt"
	"os"
	"bufio"
	"net"
	"strings"
)

func main()  {
	//var(
	//	msg string
	//	err error
	//)
	//fmt.Printf("请输入消息:\n")
	//for {
	//
	//	if msg, err = bufio.NewReader(os.Stdin).ReadString('\n'); err != nil {
	//		fmt.Println("input error", err.Error())
	//		break
	//	}
	//	fmt.Println("msg:",msg)
	//}
	//input := bufio.NewScanner(os.Stdin)
	//fmt.Printf("input:\r\n")
	//for input.Scan(){
	//	fmt.Println("hello:",input.Text())
	//}

	 connServer()
}

func connServer()  {
	var(
		tcpAddr *net.TCPAddr
		tcpConn *net.TCPConn
		err error
	)
	const TCPADDR  = "127.0.0.1:2000"
     if tcpAddr,err = net.ResolveTCPAddr("tcp",TCPADDR);err!=nil{
     	fmt.Println("connect err ",err.Error())
     	return
	 }
	 if tcpConn,err = net.DialTCP("tcp",nil,tcpAddr);err!=nil{
		 fmt.Println("DialTCP err",err.Error())
		 return
	 }
	  sendMsg(tcpConn)
}

func sendMsg(conn *net.TCPConn)  {
	var(
		msg string
		reMsg string
		err error
	)
	fmt.Printf("请输入消息:\n")
	for{
		if msg, err= bufio.NewReader(os.Stdin).ReadString('\n');err!=nil{
			fmt.Println("input error",err.Error())
			break
		}
		if _,err = conn.Write([]byte(msg));err!=nil{
			 fmt.Println("error:",err.Error())
			 goto STARTFOR
		}
		if strings.EqualFold(strings.TrimSpace(msg),"bye"){
			conn.Close()
			break
		}
		//line, _, _ := bufio.NewReader(conn).ReadLine()
		//line,_:=bufio.NewReader(conn).ReadBytes('\n')
		//bufio.NewReader(conn).Read()
		if reMsg,err = bufio.NewReader(conn).ReadString('\n');err!=nil{
             fmt.Println("readstring error:",err.Error())
			 goto STARTFOR
		}
		fmt.Println(reMsg)

		STARTFOR:
			fmt.Printf("请输入消息:\n")
	}
	os.Exit(1)
}

