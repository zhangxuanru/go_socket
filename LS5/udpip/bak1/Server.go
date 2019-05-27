package main

import (
	"time"
	"net"
	"socket/LS5/lib/protocol"
	"fmt"
	"strings"
)

func main()  {
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
		go handleTcpClient(tcpConn)
	}
}

func handleTcpClient(conn net.Conn)  {
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




