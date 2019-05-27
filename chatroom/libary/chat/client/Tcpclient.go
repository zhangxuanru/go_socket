package client

import (
	"net"
	"socket/chatroom/conf"
	"fmt"
	"os"
	"bufio"
	"strings"
	"os/signal"
	"syscall"
	"socket/chatroom/libary/protocol"
)

type TcpClient struct {
    serverTcpAddr *net.TCPAddr
    serverTcpConn *net.TCPConn
    inputMsg chan string
    readServerMsg chan string
    signal chan os.Signal
    closeChan chan bool
	closeStatus bool
}

var(
	clientStatus map[string]bool
)
func InitTcpClient()  {
	tcp:= TcpClient{
		inputMsg:make(chan string),
		readServerMsg:make(chan string),
		closeChan:make(chan bool),
		signal: make(chan os.Signal),
		closeStatus:false,
	}
	clientStatus = make(map[string]bool)
	signal.Notify(tcp.signal, syscall.SIGINT, syscall.SIGTERM)
	tcp.connServer()
	tcp.handle()
}

var(
	err error
)

func (client *TcpClient) connServer()  {
	var(
		clientKey string
	)
	if client.serverTcpAddr,err = net.ResolveTCPAddr("tcp",conf.GetTcpServerAddress());err!=nil{
          goto printErr
	}
	if client.serverTcpConn,err = net.DialTCP("tcp",nil,client.serverTcpAddr);err!=nil{
		goto printErr
	}
	clientKey = client.genKey()
	clientStatus[clientKey] = true
	fmt.Printf("连接 %s 成功\n",conf.GetTcpServerAddress())
    return
	printErr:
		fmt.Println("error:",err.Error())
 	    return
}

func (client *TcpClient) handle()  {
	var(
		msg string
	)
	if client.serverTcpConn == nil{
		fmt.Println("error  serverTcpConn nil")
		return
	}
     go client.readMsg()
     go client.intoMsgChan()
     for{
		 select {
		 case msg = <- client.inputMsg:
		 	     client.writeMsg(msg)
		 case msg = <-client.readServerMsg:
			    fmt.Printf("[%s:]:%s\n",client.serverTcpConn.RemoteAddr().String(),msg)
		 case <-client.signal:
			     client.writeMsg("bye")
			     goto close
		 case <-client.closeChan:
		 	    fmt.Println("close chan")
			    goto close
		 }
	 }
	 return
	close:
		client.close()
		fmt.Println("hand close")
}

func (client *TcpClient) writeMsg(msg string)  {
	   data := protocol.Pack([]byte(msg))
       client.serverTcpConn.Write(data)
}

func (client *TcpClient) readMsg()  {
	var(
		buf []byte
		tmpBuf []byte
		readChan chan string
	)
	err = nil
	clientKey := client.genKey()
	buf = make([]byte,1024)
	tmpBuf = make([]byte,0)
	readChan = make(chan string)
	go client.readServerRecvMsg(readChan)
    for{
		if client.serverTcpConn == nil || clientStatus[clientKey] == false{
			goto close
		}
    	   if client.serverTcpConn == nil{
    		  goto close
		   }
      	   if _,err = client.serverTcpConn.Read(buf);err!=nil{
      	   	   goto close
		   }
	  	   tmpBuf = protocol.Unpack(append(tmpBuf,buf...),readChan)
	  }
	close:
		client.close()
	    client.sendCloseChan()
		fmt.Println("read close")
}

func (client *TcpClient) readServerRecvMsg(readChan chan string)  {
	var(
		msg string
	)
	for{
		select {
		case msg = <-readChan:
			fmt.Println("msg",msg)
			client.readServerMsg<-msg
			if strings.EqualFold(msg,"bye"){
				goto close
			}
		}
	}
close:
	client.close()
	client.sendCloseChan()
	fmt.Println("read close")
}

func (client *TcpClient) intoMsgChan()  {
	 var(
	 	msg string
	 )
	clientKey := client.genKey()
	 for{
	 	if client.serverTcpConn == nil || clientStatus[clientKey] == false {
			goto close
		}
		 if msg,err = getStdinMsg();err!=nil{
               goto close
		 }
		 msg = strings.TrimSpace(msg)
		 client.inputMsg <-msg
		 if strings.EqualFold(msg,"bye"){
             //goto close
		 }
	 }
	close:
		client.sendCloseChan()
		client.close()
	    fmt.Println("input close")
}

func (client *TcpClient) close()  {
     if client.serverTcpConn!=nil{
		   clientKey := client.genKey()
		   clientStatus[clientKey] = false
     	   client.serverTcpConn.Close()
	 }
}

func (client *TcpClient) sendCloseChan()  {
	  clientKey := client.genKey()
	if clientStatus[clientKey] == false{
		 client.closeChan<-true
	}
	fmt.Println("send close chan",clientStatus[clientKey])
}


func getStdinMsg() (msg string ,err error) {
	return  bufio.NewReader(os.Stdin).ReadString('\n')
}

func (client *TcpClient) genKey() (string) {
    return client.serverTcpConn.LocalAddr().String()
}


