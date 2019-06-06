package tcp

import (
	"fmt"
	"net"
	"socket/app/common"
	"socket/app/protocol/io"
	"socket/app/service"
)

var err error

type Client struct {
	service.SocketBase
	service.ClientMsg
}

func Init() {
	clientMsg := service.NewClientMsg()
	c := &Client{
		ClientMsg: clientMsg,
	}
	if c.connServer() != nil {
		fmt.Println("连接TCP服务器失败")
		return
	}
	c.handle()
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
	stdInIo = io.NewStdInIo()
	socketIo := service.NewSocketIo(c.SocketBase)
	go stdInIo.OutStdInMsgByChan(c.InputMsgChan)
	go socketIo.Read()

	for {
		select {
		case c.ClientMsg.InputMsg = <-c.InputMsgChan:
			socketIo.WriteData(c.ClientMsg.InputMsg)
		}
	}

}
