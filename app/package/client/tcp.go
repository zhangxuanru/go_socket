package client

import (
	"fmt"
	"net"
	"socket/app/common"
	"socket/app/package/io"
	"socket/app/package/socket"
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
	var (
		stdInIo *io.StdInIo
	)
	if c.TcpConn == nil {
		fmt.Println("连接出错.....")
		return
	}
	stdInIo = io.NewStdInIo(true)
	socketIo := socket.NewSocketIo(c.SocketBase, socket.SEND_STRING)
	go stdInIo.OutStdInMsgByChan(c.InputMsgChan)
	go socketIo.SocketPack.Read(c.TcpConn, c.ReadMsgChan)

	for {
		select {
		case c.ClientMsg.InputMsg = <-c.InputMsgChan:
			socketIo.WriteData(c.ClientMsg.InputMsg)
			common.CheckBye([]byte(c.ClientMsg.InputMsg), c.IsCloseChan)
		case c.ClientMsg.ReadMsg = <-c.ReadMsgChan:
			if socketIo.IsStrByCurrPack(c.ClientMsg.ReadMsg) == false {
				socketIo.ResetSocketPack()
			}
			c.ClientMsg.ReadMsg = common.RemoveStrSendHeader(c.ClientMsg.ReadMsg)
			socketIo.SocketPack.Receive(c.ClientMsg.ReadMsg)
			common.CheckBye(c.ClientMsg.ReadMsg, c.IsCloseChan)
		case <-c.IsCloseChan:
			goto END
		}
	}
END:
	fmt.Println("handle end")
	stdInIo.Close()
	socketIo.SocketPack.Close(c.TcpConn)
}
