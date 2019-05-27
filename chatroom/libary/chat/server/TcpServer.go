package server

import (
	"net"
	"socket/chatroom/conf"
	"fmt"
	"strings"
	"bufio"
	"os"
	"socket/chatroom/libary/protocol"
)

type tcp struct {
     tcpAddr *net.TCPAddr
     tcpListen *net.TCPListener
     conn *Conn
}

type Conn struct {
	 clientInfo chan clientInfo
	 clientMsg chan string
	 inputMsg  chan string
	 retMsg    string
}

type clientInfo struct {
	conn *net.TCPConn
	clientMsg string
}

var(
	G_tcp *tcp
	err error
	clients map[string]*net.TCPConn
)

func InitTcp()  {
	 G_tcp = &tcp{
	 	conn:&Conn{
		   clientInfo :make(chan clientInfo),
           clientMsg:make(chan string),
           inputMsg:make(chan string),
		},
	 }
	clients = make(map[string]*net.TCPConn)
	G_tcp.Listen()
}

func (t *tcp)  Listen() {
	var(
		conn *net.TCPConn
	)
	if t.tcpAddr,err =  net.ResolveTCPAddr("tcp",conf.GetTcpServerAddress());err!=nil{
          goto PrintErr
	}
	if t.tcpListen,err = net.ListenTCP("tcp",t.tcpAddr);err!=nil{
		  goto PrintErr
	}
	fmt.Println("服务器信息,Network:",t.tcpListen.Addr().Network(),"  String:",t.tcpListen.Addr().String())
	defer t.tcpListen.Close()

	go t.getWriteMsg()
	for{
		if conn,err = t.tcpListen.AcceptTCP();err!=nil{
             continue
		}
		AddClient(conn)
		go t.handleClient(conn)
	}
	return
	PrintErr:
		fmt.Println("error:",err.Error())
	    return
}


func (t *tcp) handleClient(conn *net.TCPConn)  {
	fmt.Println("Server conn RemoteAddr:",conn.RemoteAddr().String(),"  LocalAddr:",conn.LocalAddr().String())
    var(
    	msg string
    	clientInfo clientInfo
	)
	go t.readMsg(conn)
	for{
		select {
		case clientInfo = <-t.conn.clientInfo:
			   msg = clientInfo.clientMsg
			   printMsgInfo(clientInfo.clientMsg,clientInfo.conn)
			   t.WriteChatRoomMsg(clientInfo.conn,msg) //群发
		       if checkIsBye(msg){
				   goto Start
			   }
		case msg = <-t.conn.inputMsg:
			    //t.WriteMsg(conn,msg)  //单发
			     t.WriteRoomMsg(conn,msg)
				if checkIsBye(msg){
					goto Start
				}
		}
	}
	return
	Start:
		fmt.Println("bye bye:")
	    return
}

func (t *tcp) readMsg(conn *net.TCPConn)  {
	var(
		buf []byte

		err error
		tmpBuf []byte
		readChan chan string
	)
	buf = make([]byte,1024)
	tmpBuf = make([]byte,0)
	readChan = make(chan string)
	go t.readChanMsg(readChan,conn)
	for{
    	if conn == nil {
    		goto ConnClose
		}
        //bufio.NewReader(conn).ReadString('\n')
		//bufio.NewReader(conn).Read(buf) //读取指定字节数
         if _,err = conn.Read(buf);err!=nil{
         	 continue
		 }
		 fmt.Println("buf:",string(buf))
		 tmpBuf = protocol.Unpack(append(tmpBuf,buf...),readChan)
	   }
	 return
	ConnClose:
		fmt.Println("read close chan ")
		t.Close(conn)
	    return
}


func (t *tcp) readChanMsg(readChan chan string,conn *net.TCPConn)  {
	var(
		readMsg string
	)
	  for{
		  select {
		  case readMsg = <-readChan:
			  t.conn.clientInfo <- clientInfo{
				  clientMsg:readMsg,
				  conn:conn,
			  }
			  if checkIsBye(readMsg){
				  goto ConnClose
			  }
		  }
	  }
 ConnClose:
	fmt.Println("read close chan ")
	t.Close(conn)
	return
}

func (t *tcp) WriteMsg(conn *net.TCPConn,msg string)  {
	conn.Write([]byte(msg))
}

func (t *tcp) WriteChatRoomMsg(conn *net.TCPConn,msg string)  {
		if strings.EqualFold(msg,"bye") {
			  return
		}
        for _,client := range clients{
			if checkClientEqual(client,conn) == false {
				   fmt.Println(client.RemoteAddr().String())
				   sendMsg := protocol.Pack([]byte(msg))
        	       client.Write(sendMsg)
			}
		}
}

func (t *tcp) WriteRoomMsg(conn *net.TCPConn,msg string)  {
	for _,client := range clients{
		    sendMsg := protocol.Pack([]byte(msg))
			client.Write(sendMsg)
			if strings.EqualFold(msg,"bye") {
				t.Close(client)
			}
	}
}

func (t *tcp) getWriteMsg()  {
	var(
		msg string
	)
	for{
	    if msg,err = getStdinMsg();err!=nil{
             goto Loop
		 }
		msg = strings.TrimSpace(msg)
		t.conn.inputMsg <- msg
	}
	Loop:
		fmt.Println("input bye")
	    return
}


func (t *tcp) Close(conn *net.TCPConn)  {
	if 	conn!=nil{
	      DelClient(conn)
		  err = conn.Close()
		  fmt.Println("close:",conn.RemoteAddr().String(),"close error:",err)
	}
	conn = nil
	fmt.Println("close conn")
}


func checkIsBye(msg string) bool {
	return strings.EqualFold(msg,"bye")
}

func getStdinMsg() (msg string ,err error) {
	 return  bufio.NewReader(os.Stdin).ReadString('\n')
}

func printMsgInfo(msg string,conn *net.TCPConn)  {
	fmt.Printf("[%s]:%s\n",conn.RemoteAddr().String(),msg)
}

func AddClient(conn *net.TCPConn)  {
	 key := generyTcpClientKey(conn)
	 if clientsHasConn(conn) == false{
		 clients[key] = conn
	 }
}

func DelClient(conn *net.TCPConn)  {
	if conn == nil{
		return
	}
	key := generyTcpClientKey(conn)
	if clientsHasConn(conn) == true{
		delete(clients,key)
	}
}

func generyTcpClientKey(conn *net.TCPConn) (key string) {
	return conn.RemoteAddr().String()
}

func clientsHasConn(conn *net.TCPConn) (bool) {
	key := generyTcpClientKey(conn)
	if _,ok:= clients[key];ok == true{
		return true
	}
	return false
}

func checkClientEqual(addr1,addr2 *net.TCPConn) bool {
	   addr1Key := generyTcpClientKey(addr1)
	   addr2Key := generyTcpClientKey(addr2)
	   return  strings.EqualFold(addr1Key,addr2Key)
}
