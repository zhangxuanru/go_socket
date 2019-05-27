package server

import (
	"net"
	"socket/chatroom/conf"
	"fmt"
	"socket/chatroom/libary/util"
	"strings"
	"source/openbilibili-go-common/app/service/ops/log-agent/pkg/bufio"
	"os"
	"time"
)

type UdpServer struct {
	udpAddr *net.UDPAddr
	udpConn *net.UDPConn
	readMsgChan chan ReadMsg
	WriteMsgChan chan string
	ClientList   map[string]ClientInfo
}

type ReadMsg struct {
	msg string
	udpAddr *net.UDPAddr
}

type ClientInfo struct {
	udpAddr *net.UDPAddr
	connectTime int64
}

var ClientList map[string]ClientInfo


var G_UdpServer *UdpServer


func InitUdpServer()  {
	var(
		err error
	)
	G_UdpServer = &UdpServer{
		readMsgChan: make(chan ReadMsg),
		WriteMsgChan:make(chan string),
	}
	ClientList = make(map[string]ClientInfo)

	if err = G_UdpServer.Start();err!=nil{
		goto PRINTERR
	}
	G_UdpServer.Listen()
	PRINTERR:
		fmt.Println("error:",err.Error())
	    return
}


func (udp *UdpServer) Start() (err error) {
	  if udp.udpAddr,err =  net.ResolveUDPAddr("udp",conf.GetUdpServerAddress());err!=nil{
           return
	  }
	  udp.udpConn,err = net.ListenUDP("udp",udp.udpAddr)
	  return
}

func (udp *UdpServer) Listen()  {
	fmt.Println("UDP服务端启动成功 ip:",udp.udpAddr.IP.String(), " Port:",udp.udpAddr.Port)
	defer udp.udpConn.Close()
	for{
		  udp.handleClient()
	  }
}


func (udp *UdpServer) handleClient()  {
	var(
		receivedMsg ReadMsg
		inputMsg string
	)
    go udp.ReadClientMsg()
    go udp.GetInputMsg()
	select {
	case receivedMsg= <-udp.readMsgChan:
             udp.PrintClientMsg(receivedMsg)
	         udp.WriteRootMsg(receivedMsg)
	case inputMsg = <-udp.WriteMsgChan:
		     udp.WriteBroadcastMsg(inputMsg)
	}
}


func (udp *UdpServer) ReadClientMsg()  {
	var(
		buf []byte
		msg string
		udpAddr *net.UDPAddr
		err error
	)
	buf = make([]byte,1024)
	if _,udpAddr,err = udp.udpConn.ReadFromUDP(buf);err!=nil{
		fmt.Println("ReadFromUDP error:",err.Error())
		return
	}
	msg = util.ByteToString(buf)
	udp.readMsgChan <- ReadMsg{
		msg:msg,
		udpAddr:udpAddr,
	}
	udp.AddClient(udpAddr)
	if strings.EqualFold(msg,"bye"){
         udp.Close(udpAddr)
	}
	//fmt.Println("count client:",len(ClientList))
}


func (udp *UdpServer) GetInputMsg()  {
	var(
		msg string
		err error
	)
	if msg,err = getStdinValue();err!=nil{
		return
	}
	udp.WriteMsgChan <- msg
}


func (udp *UdpServer) WriteBroadcastMsg(msg string)  {
	for _,Client := range ClientList{
		udp.udpConn.WriteToUDP([]byte(msg),Client.udpAddr)
	}
}


func (udp *UdpServer) WriteRootMsg(receivedMsg ReadMsg)  {
	for _,Client := range ClientList{
		   if udp.CheckClientEqual(receivedMsg.udpAddr,Client.udpAddr){
		   	   continue
		   }
		receivedMsg.msg = fmt.Sprintf("%s[%s]",receivedMsg.msg,receivedMsg.udpAddr.String())
		udp.udpConn.WriteToUDP([]byte(receivedMsg.msg),Client.udpAddr)
	}
}



func (udp *UdpServer) AddClient(udpAddr *net.UDPAddr)  {
	 key := generyClientKey(udpAddr)
	 if _,ok:= ClientList[key];ok==true{
            return
	 }
	ClientList[key] = ClientInfo{
		udpAddr:udpAddr,
		connectTime:time.Now().Unix(),
	}
}


func (udp *UdpServer) RemoveClient(addr *net.UDPAddr)  {
	key := generyClientKey(addr)
	if _,ok:= ClientList[key];ok==true{
		 delete(ClientList,key)
	}
}

func (udp *UdpServer) CheckClientExists(addr *net.UDPAddr) bool {
	key := generyClientKey(addr)
	_,ok:= ClientList[key]
	return  ok
}

func (udp *UdpServer) CheckClientEqual(addr1,addr2 *net.UDPAddr) bool {
	addr1Key := generyClientKey(addr1)
	addr2Key := generyClientKey(addr2)
    return  strings.EqualFold(addr1Key,addr2Key)
}



func (udp *UdpServer) Close(addr *net.UDPAddr)  {
	udp.RemoveClient(addr)
    fmt.Println("close:",addr.IP.String(),"port:",addr.Port)
}


func (udp *UdpServer) PrintClientMsg(msg ReadMsg)  {
	fmt.Printf("[%s:%d]:%s\n",msg.udpAddr.IP.String(),msg.udpAddr.Port, msg.msg)
}


func getStdinValue() (msg string,err error) {
	//fmt.Printf("my:")
	msg,err = bufio.NewReader(os.Stdin).ReadString('\n')
	return
}


func generyClientKey(udpAddr *net.UDPAddr) (key string) {
	key = fmt.Sprintf("%s:%d",udpAddr.IP.String(),udpAddr.Port)
	return
}
