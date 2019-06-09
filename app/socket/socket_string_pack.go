package socket

import (
	"errors"
	"net"
	"fmt"
)

type stringPack struct {
}

func (s *stringPack) Write(conn *net.TCPConn, data []byte) (int, error) {
	if conn == nil {
		return 0, errors.New("conn is nil")
	}
	return conn.Write(data)
}

func (s *stringPack) Read(conn *net.TCPConn,readChan chan []byte) (err error) {
	if conn == nil {
		return  errors.New("conn is nil")
	}
	for{
		if conn == nil{
			 goto CLOSE
		}
		buf := make([]byte, 1024)
		if _, err = conn.Read(buf); err != nil {
			 goto CLOSE
		}

		readChan <- buf
	}
	return   err
	CLOSE:
		fmt.Println("read close error:",err)
	    //readChan<-[]byte("bye")
	return  err
}

func (s *stringPack) Close(conn *net.TCPConn) error {
	return conn.Close()
}


func (s *stringPack) Receive(receiveMsg []byte)  {
 fmt.Printf("receiveMsg:%s\n",string(receiveMsg))
}
