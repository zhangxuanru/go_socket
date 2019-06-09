package client

import (
	"fmt"
	"net"
	"socket/app/common"
	"socket/app/protocol/io"
	"socket/app/socket"
)

var err error

type Client struct {
	 socket.SocketBase
	*socket.ClientMsg
}

func Init() {
	c := &Client{
		ClientMsg: socket.NewClientMsg(),
	}
	if c.connServer() != nil {
		fmt.Println("连接TCP服务器失败")
		return
	}
	c.handle()
	fmt.Println("bye bye")
}

func (c *Client) connServer() error {
	if c.TcpAddr, err = net.ResolveTCPAddr(common.GetNetWorkType(), common.GetServerAddress()); err != nil {
		goto PrintErr
	}
	if c.TcpConn, err = net.DialTCP(common.GetNetWorkType(), nil, c.TcpAddr); err != nil {
		goto PrintErr
	}
	fmt.Println("连接成功:", common.GetServerAddress())
	return nil
PrintErr:
	common.PrintErr(err)
	return err
}

func (c *Client) handle() {
	fmt.Println("handle........")
	var (
		stdInIo *io.StdInIo
	)
	if c.TcpConn == nil {
		fmt.Println("连接出错.....")
		return
	}
	stdInIo = io.NewStdInIo()
	socketIo := socket.NewSocketIo(c.SocketBase)
	go stdInIo.OutStdInMsgByChan(c.InputMsgChan)
	go socketIo.SocketPack.Read(c.TcpConn,c.ReadMsgChan)

	for {
		select {
		case c.ClientMsg.InputMsg = <-c.InputMsgChan:
			  socketIo.WriteData(c.ClientMsg.InputMsg)
		      socketIo.CheckBye([]byte(c.ClientMsg.InputMsg),c.IsCloseChan)
		case c.ClientMsg.ReadMsg = <-c.ReadMsgChan:
			  socketIo.SocketPack.Receive(c.ClientMsg.ReadMsg)
			  socketIo.CheckBye(c.ClientMsg.ReadMsg,c.IsCloseChan)
		case  <-c.IsCloseChan:
               goto  END
		}
	}
	END:
		fmt.Println("handle end")
	    stdInIo.Close()
	    socketIo.SocketPack.Close(c.TcpConn)
}
