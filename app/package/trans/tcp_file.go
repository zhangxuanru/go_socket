package trans

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"socket/app/common"
	"socket/app/config"
)

type TcpFilePack struct {
	closeChan chan bool
}

func (s *TcpFilePack) Write(conn *net.TCPConn, data []byte) (int, error) {
	if conn == nil {
		return 0, errors.New("conn is nil")
	}
	SendFile(conn, data)

	fmt.Println(string(data))

	data = append([]byte(config.SEND_FILE_HEADER_PACK), data...)
	return conn.Write(Pack(data))
}

func (s *TcpFilePack) Read(conn *net.TCPConn, readMsgChan chan []byte) (err error) {
	NewRead().Read(conn, readMsgChan)
	return nil
}

func (s *TcpFilePack) Close(conn *net.TCPConn) error {
	fmt.Println("conn is close")
	return conn.Close()
}

func (s *TcpFilePack) Receive(receiveMsg []byte) {
	fmt.Println("file:----")
	fmt.Println("receiveMsg:", string(receiveMsg))
}

func SendFile(conn *net.TCPConn, fileMsg []byte) (err error) {
	var (
		fileInfo *os.File
	)
	fileMsg = common.RemoveFilePrefixMsg(fileMsg)
	filePath := string(fileMsg)
	if common.IsFileNotExists(filePath) == true {
		err = errors.New("file is not exists")
		return
	}
	if fileInfo, err = os.Open(filePath); err != nil {
		return
	}
	_, file := filepath.Split(filePath)
	defer fileInfo.Close()
	buf := make([]byte, 1024)
	for {
		n, e := fileInfo.Read(buf)
		if e == io.EOF {
			break
		}
		conn.Write(buf[:n])
	}
	conn.Write([]byte("end:" + file))
	return
}

func NewTcpFilePack() *TcpFilePack {
	return &TcpFilePack{
		closeChan: make(chan bool),
	}
}
