package trans

import (
	"errors"
	"fmt"
	"net"
)

type TcpRead struct {
	closeChan chan bool
}

func (r *TcpRead) Read(conn *net.TCPConn, readChan chan []byte) (err error) {
	var (
		tmpBuf []byte
		buf    []byte
	)
	if conn == nil {
		return errors.New("conn is nil")
	}
	tmpBuf = make([]byte, 0)
	for {
		if conn == nil {
			goto CLOSE
		}
		select {
		case <-r.closeChan:
			fmt.Println("---closeChan--")
			err = errors.New("auto close")
			goto CLOSE
		default:
			buf = make([]byte, 1024)
			if _, err = conn.Read(buf); err != nil {
				goto CLOSE
			}
			//解包
			tmpBuf = UnPack(append(tmpBuf, buf...), readChan)
		}
	}
	return err
CLOSE:
	tmpBuf = make([]byte, 0)
	buf = make([]byte, 0)
	fmt.Println("read close......")
	return err
}

func (r *TcpRead) Close() {
	go func() {
		r.closeChan <- true
	}()
}

func NewRead() *TcpRead {
	return &TcpRead{
		closeChan: make(chan bool),
	}
}
