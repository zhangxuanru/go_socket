package client

import (
	"net"
	"socket/chatroom/conf"
	"fmt"
	"socket/chatroom/libary/util"
	"bufio"
	"os"
	"strings"
	"syscall"
	"os/signal"
)

type ClientMsg struct {
	inputMsg string
}

type ServerMsgInfo struct {
	 udpAddr *net.UDPAddr
	 receivedMsg  string
}

type Client struct {
    udpAddr *net.UDPAddr
    udpConn *net.UDPConn
	clientMsgInfo chan ClientMsg
	serverMsgInfo  chan ServerMsgInfo
	signalChan     chan os.Signal
	done         bool
}

var G_client *Client

func InitClient()  {
	G_client = &Client{
		clientMsgInfo: make(chan ClientMsg),
		serverMsgInfo:make(chan ServerMsgInfo),
		signalChan: make(chan os.Signal),
		done:false,
	}
	signal.Notify(G_client.signalChan, syscall.SIGINT, syscall.SIGTERM)
	G_client.connServer()
	if G_client.udpConn != nil{
		 G_client.handleServer()
	}
}

func (client *Client) connServer()  {
	var(
		udpAddr *net.UDPAddr
		udpConn *net.UDPConn
		err error
	)
	if udpAddr,err = net.ResolveUDPAddr("udp",conf.GetUdpServerAddress());err!=nil{
		fmt.Println("connServer error:",err.Error())
		return
	}
	if udpConn,err = net.DialUDP("udp",nil,udpAddr);err!=nil{
		fmt.Println("DialUDP error:",err.Error())
		return
	}
	fmt.Println("connect udp server:",udpConn.RemoteAddr().String())
	client.udpAddr = udpAddr
	client.udpConn = udpConn
}


func  (client *Client) handleServer()  {
	var(
		serverMsgInfo ServerMsgInfo
		clientMsgInfo  ClientMsg
	)
      go client.readServerMsg()
      go client.inputMsg()

      for{
      	  if client.done == true{
      	  	  goto CloseConn
		  }
		  select {
		     case serverMsgInfo = <-client.serverMsgInfo:
		     	  client.PrintReceivedMsg(serverMsgInfo.receivedMsg)
		     case clientMsgInfo = <-client.clientMsgInfo:
		     	   client.sendMsg(clientMsgInfo)
		     case <-client.signalChan:
					 clientMsgInfo = ClientMsg{
						inputMsg:"bye",
					  }
				     client.sendMsg(clientMsgInfo)
		     	     goto CloseConn
		}
	  }
	CloseConn:
		client.udpConn.Close()
	    fmt.Println("bye bye")
}


func (client *Client) readServerMsg()  {
	var(
		err error
		receivedMsg string
		udpServerAddr *net.UDPAddr
	)
	for{
		if client.done == true{
			return
		}
		var(
			buf []byte
		)
		 buf = make([]byte,1024)
		 if _,udpServerAddr,err = client.udpConn.ReadFromUDP(buf);err!=nil{
			   goto PRINTERR
		 }
		 receivedMsg = util.ByteToString(buf)
		 client.serverMsgInfo <- ServerMsgInfo{
			udpAddr:udpServerAddr,
			receivedMsg:receivedMsg,
		 }
	}
	 PRINTERR:
	 	fmt.Println("error:",err.Error())
	    return
}


func (client *Client) sendMsg(msg ClientMsg)  {
	var(
		err error
	)
    if _,err = client.udpConn.Write([]byte(msg.inputMsg));err!=nil{
    	fmt.Println("error:",err.Error())
	}
}


func (client *Client) inputMsg()  {
	var(
		msg string
		err error
	)
	for{
		if msg,err = bufio.NewReader(os.Stdin).ReadString('\n');err!=nil{
			fmt.Println("error:",err.Error())
			return
		}
		msg = strings.TrimSpace(msg)
		if len(msg) == 0{
			continue
		}
		client.clientMsgInfo<- ClientMsg{
			 inputMsg:msg,
		}
		if strings.EqualFold(msg,"bye"){
			 goto closeInput
		}
	}
	closeInput:
		 fmt.Println("close input")
	     client.done = true
}

func (client *Client) PrintReceivedMsg(msg string ) {
	ipPos := strings.Index(msg,"[")
	if ipPos > -1{
		fmt.Printf("%s:%s\n",msg[ipPos:],msg[0:ipPos])
	}else{
        fmt.Printf("[server:]%s\n",msg)
	}
}



