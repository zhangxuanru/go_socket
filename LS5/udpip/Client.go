package main

import (
	"net"
	"socket/LS5/lib/protocol"
	"fmt"
	"socket/LS5/conf"
	"bytes"
	"time"
	"errors"
	"bufio"
	"os"
	"strings"
)

type TcpServer struct {
	inputMsgChan chan string
	readMsgChan chan string
	closeConn bool
}
var(
	G_TcpServerClient *TcpServer
)

func main()  {
	run()
}

func run()  {
	var(
		tcpAddr net.TCPAddr
		err error
	)
	if tcpAddr,err = searchUdpServer();err!=nil{
		fmt.Println("searchUdpServer error:",err.Error())
		return
	}
	connectTcpServer(&tcpAddr)
}

func searchUdpServer() (tcpServer net.TCPAddr,err error) {
	var(
		remoteUdpAddr *net.UDPAddr
		localUdpAddr *net.UDPAddr
		udpConn *net.UDPConn
		clientAddr string
		tcpServerChan chan net.TCPAddr
	)
	tcpServerChan = make(chan net.TCPAddr)
	if remoteUdpAddr,err = net.ResolveUDPAddr("udp",protocol.GetUdpBroadcastAddress());err!=nil{
		fmt.Println("ResolveUDPAddr remoteUdpAddr error:",err.Error())
		return
	}
	clientAddr = fmt.Sprintf(":%d",conf.UDP_CLIENT_PORT)
	if localUdpAddr,err = net.ResolveUDPAddr("udp",clientAddr);err!=nil{
		fmt.Println("ResolveUDPAddr localUdpAddr error:",err.Error())
		return
	}
	if udpConn,err = net.ListenUDP("udp",localUdpAddr);err!=nil{
		fmt.Println("DialUDP error:",err.Error())
		return
	}
	defer udpConn.Close()

	go  sendMsgServer(udpConn,remoteUdpAddr,tcpServerChan)
	return getTcpServer(tcpServerChan)
}

func getTcpServer(tcpServerChan chan net.TCPAddr) (tcpServer net.TCPAddr, err error) {
	var(
		ret chan bool
		timer *time.Timer
	)
	timer = time.NewTimer(time.Second * conf.UDP_SEARCH_TIMEOUT)
	ret = make(chan bool)
	go func(tcpServerChan chan net.TCPAddr) {
		select {
		case <-timer.C:
			err = errors.New("search udp time out")
			ret<-true
		case tcpServer = <-tcpServerChan:
			ret<-true
		}
	}(tcpServerChan)
	<-ret
	return tcpServer,err
}

func sendMsgServer(conn *net.UDPConn,addr *net.UDPAddr,tcpServerChan chan net.TCPAddr) (TCPAddr net.TCPAddr) {
	var(
		buf []byte
		serverUdpAddr *net.UDPAddr
		err error
		receiveMsg string
		sendMsg string
		msgHeaderLen int
	)
	buf = make([]byte,1024)
	msgHeaderLen = len(protocol.UDP_MSG_HEADER)

	sendMsg = protocol.BuildWithPort(conf.UDP_CLIENT_PORT)
	buf = bytes.Replace(buf,buf[0:msgHeaderLen],protocol.UDP_MSG_HEADER,1)
	buf = bytes.Replace(buf,buf[msgHeaderLen:],[]byte(sendMsg),1)

	if _,err = conn.WriteToUDP(buf,addr);err!=nil{
		fmt.Println("WriteToUDP error:",err.Error())
	}
	buf = make([]byte,1024)
	if _,serverUdpAddr,err = conn.ReadFromUDP(buf);err!=nil{
		fmt.Println("ReadFromUDP error:",err.Error())
		return
	}
	receiveMsg = protocol.GetStringByBuF2(buf)
	fmt.Println("server ip:",serverUdpAddr.IP," server port:",serverUdpAddr.Port," server msg:",receiveMsg)

	TCPAddr = net.TCPAddr{
		IP:serverUdpAddr.IP,
		Port:conf.TCP_SERVER_PORT,
	}
	tcpServerChan <- TCPAddr
	return
}


//connect tcp server
func connectTcpServer(tcpAddr *net.TCPAddr)  {
	var(
		tcpConn *net.TCPConn
		err error
	)
	if tcpConn,err = net.DialTCP("tcp",nil,tcpAddr);err!=nil{
		fmt.Println("DialTCP error:",err.Error())
		return
	}
	InitTcpConn()
	sendTcpMsgScheduler(tcpConn)
}

func InitTcpConn()  {
	G_TcpServerClient = &TcpServer{
            inputMsgChan:make(chan string),
            readMsgChan:make(chan string),
		    closeConn:false,
	}
}


func sendTcpMsgScheduler(conn *net.TCPConn)  {
	var(
		msg string
	)
	go getStdinInputMsg()
	go getServerReturnsMsg(conn)
	for{
		if G_TcpServerClient.closeConn==true{
			break
		}
		select {
		     case msg = <-G_TcpServerClient.inputMsgChan:
				    sendTcpMsg(conn,msg)
		    case msg = <-G_TcpServerClient.readMsgChan:
		    	   fmt.Println("server:",msg)
		}
	}
	fmt.Println("close tcp conn")
	if conn!=nil{
		conn.Close()
	}
}


func getStdinInputMsg()  {
	var(
		msg string
		err error
	)
	for{
		if G_TcpServerClient.closeConn == true{
			goto START
		}
		fmt.Println("请输入要发送的消息:")
		if msg,err = bufio.NewReader(os.Stdin).ReadString('\n');err==nil{
			msg = strings.TrimSpace(msg)
			G_TcpServerClient.inputMsgChan <- msg
		}
		if strings.EqualFold(msg,"bye"){
			G_TcpServerClient.closeConn = true
			goto START
		}
	}
START:
	fmt.Println("close input")
}


func sendTcpMsg(conn *net.TCPConn,msg string)  {
	var (
		err error
	)
	if _,err = conn.Write([]byte(msg));err!=nil{
		fmt.Println("Write error:",err.Error())
	}
}


func getServerReturnsMsg(conn *net.TCPConn)  {
	var(
		buf []byte
		msg string
		err error
	)
	for{
		if G_TcpServerClient.closeConn == true{
			break
		}
		buf = make([]byte,1024)
		if _,err = conn.Read(buf);err==nil{
				msg = protocol.GetStringByBuF2(buf)
				conn.SetDeadline(time.Now().Add(2*time.Minute))
				G_TcpServerClient.readMsgChan <- msg
	  	  }else{
			  fmt.Println("Read error:",err.Error())
			  break
		 }
	}
	fmt.Println("close read tcp conn")
}

