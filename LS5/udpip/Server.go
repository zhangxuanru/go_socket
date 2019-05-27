package main

import (
	"time"
	"net"
	"socket/LS5/lib/protocol"
	"fmt"
	"strings"
	"os"
	"bufio"
	"runtime"
)

type TCPSERVER struct {
	inputMsg chan string
	readMsg  chan string
	byeMsg   chan bool
	exitConn  bool
}

var (
	G_TCPSERVER *TCPSERVER
)

func main()  {
   runtime.GOMAXPROCS(runtime.NumCPU())
  go startUDPServer()

  go startTCPServer()

  for{
      time.Sleep(10000*time.Second)
  }

}

func startUDPServer()  {
	   var(
	   	   serverAddress string
	   	   err error
	   	   udpAddr *net.UDPAddr
	   	   udpConn *net.UDPConn
	   )
	  serverAddress = protocol.GetUdpServerAddress()
      if udpAddr,err = net.ResolveUDPAddr("udp",serverAddress);err!=nil{
      	  fmt.Println("ResolveUDPAddr error:",err.Error())
          return
	  }
	  if udpConn,err = net.ListenUDP("udp",udpAddr);err!=nil{
	  	  fmt.Println("ListenUDP error:",err.Error())
	  	  return
	  }
	  fmt.Println("UDP SERVER START ")
	  defer udpConn.Close()
	  for{
		  handleUDPClient(udpConn)
	  }
}

func handleUDPClient(conn *net.UDPConn)  {
	  var(
	  	buf []byte
	    receiveMsg string
	    retMsg string
	    clientPort int
	    udpAddr *net.UDPAddr
	    clientUdpAddr *net.UDPAddr
	  	err error
	  )
	  buf = make([]byte,1024)
	  if _,udpAddr,err = conn.ReadFromUDP(buf);err!=nil{
	  	 fmt.Println("ReadFromUDP error:",err.Error())
	  	 return
	  }
	receiveMsg = protocol.GetStringByBuF2(buf)
    fmt.Println("client ip:",udpAddr.IP," client port:",udpAddr.Port," client msg:",receiveMsg)
	if protocol.CheckReceiveMsgHeader(buf) == false{
		fmt.Println("CheckReceiveMsgHeader error")
		return
	}
	if clientPort,err = protocol.ParsePort(receiveMsg);err!=nil{
		fmt.Println("ParsePort error:",err.Error())
		return
	}
	clientUdpAddr = &net.UDPAddr{
		IP:udpAddr.IP,
		Port:clientPort,
	}
	fmt.Println(clientUdpAddr)
	retMsg = protocol.BuildWithSn("zxr----")
	conn.WriteToUDP([]byte(retMsg),clientUdpAddr)
}


func startTCPServer()  {
	var(
		tcpAddr *net.TCPAddr
		tcpListener *net.TCPListener
		tcpConn net.Conn
		err error
	)
    if tcpAddr,err = net.ResolveTCPAddr("tcp",protocol.GetTcpServerAddress());err!=nil{
    	fmt.Println("ResolveTCPAddr error:",err.Error())
    	return
	}
	if tcpListener,err = net.ListenTCP("tcp",tcpAddr);err!=nil{
		fmt.Println("ListenTCP error:",err.Error())
		return
	}
	fmt.Println("TCP SERVER START")
	for{
		if tcpConn,err = tcpListener.Accept();err!=nil{
           fmt.Println( "Accept error:",err.Error())
           continue
		}
		initTcpClient()
		go handleTcpClient(tcpConn)
	}
}

func initTcpClient()  {
	G_TCPSERVER = &TCPSERVER{
		inputMsg:make(chan string),
		readMsg:make(chan string),
		byeMsg:make(chan bool),
		exitConn:false,
	}
}

func handleTcpClient(conn net.Conn)  {
	var(
		msg string
	)
	defer conn.Close()
	go getReadMsg(conn)
	go getInputMsg()
	fmt.Println("TCP Server conn RemoteAddr:",conn.RemoteAddr().String(),"  LocalAddr:",conn.LocalAddr().String())
	for{
		select {
		case msg = <-G_TCPSERVER.readMsg:
			     fmt.Println("client msg:",msg)
		case msg = <-G_TCPSERVER.inputMsg:
			    handleWriteTcpClient(conn,msg)
		case  <-G_TCPSERVER.byeMsg:
		         goto BYE
		 }
	}
	BYE:
		fmt.Println("bye bye")
	    G_TCPSERVER.exitConn = true
}

func getInputMsg() {
	var(
		msg string
		err error
	)
	for{
		if G_TCPSERVER.exitConn == false{
				fmt.Println("请输入要回复的消息:")
				if msg,err =  bufio.NewReader(os.Stdin).ReadString('\n');err!=nil{
					return
				}
				G_TCPSERVER.inputMsg<-msg
		 }
		if  G_TCPSERVER.exitConn == true{
			  goto START
		}
	}
START:
	fmt.Println("exit input")
}


func getReadMsg(conn net.Conn) {
	var(
		buf []byte
		msg string
		err error
	)
	for{
		if  G_TCPSERVER.exitConn == true{
			goto START
		}
		if conn == nil{
			fmt.Println("conn is close ")
			break
		}
		buf = make([]byte,1024)
		if _,err = conn.Read(buf);err==nil{
			msg = protocol.GetStringByBuF2(buf)
			conn.SetReadDeadline(time.Now().Add(time.Minute * 2))
			G_TCPSERVER.readMsg <- msg
		}
		if strings.EqualFold(msg,"bye"){
			  G_TCPSERVER.byeMsg<-true
		}
	}
	START:
		fmt.Println("exit read")
}

func handleWriteTcpClient(conn net.Conn,msg string)  {
	var(
		err error
	)
	if _,err = conn.Write([]byte(msg));err!=nil{
		fmt.Println("write error:",err.Error())
	}
}



func handleTcpClient_bak(conn net.Conn)  {
	var(
		buf []byte
		receiveMsg string
		retMsg string
		err error
	)
	defer conn.Close()
	fmt.Println("TCP Server conn RemoteAddr:",conn.RemoteAddr().String(),"  LocalAddr:",conn.LocalAddr().String())
	for{
		buf = make([]byte,1024)
		if _,err = conn.Read(buf);err!=nil{
			fmt.Println("Read error:",err.Error())
			break
		}else{
			conn.SetReadDeadline(time.Now().Add(time.Minute * 2))
		}
		receiveMsg = protocol.GetStringByBuF2(buf)
		fmt.Println("tcp client msg:",receiveMsg)


		retMsg = fmt.Sprintf("tcp response:%d",len(receiveMsg))
		conn.Write([]byte(retMsg))
		if strings.EqualFold(receiveMsg,"bye"){
			fmt.Println(" close Bye")
			break
		}
	}
	fmt.Println("bye bye")
}



