package io

import (
	"bufio"
	"fmt"
	"os"
	"socket/app/common"
	"strings"
)

type StdInIo struct {
	IsCloseChan chan bool
}

func (s *StdInIo) GetStdInMsg() (string, error) {
	return bufio.NewReader(os.Stdin).ReadString('\n')
}

func (s *StdInIo) OutStdInMsgByChan(outMsg chan<- string) {
	var (
		msg string
		err error
	)
	fmt.Println("please input msg :")
	for {
		select {
		case <-s.IsCloseChan:
			goto END
		default:
			if msg, err = s.GetStdInMsg(); err != nil {
				goto END
			}
			msg = strings.TrimSpace(msg)
			outMsg <- msg
			if common.IsMsgBye(msg) {
				goto END
			}
		}
	}
END:
	fmt.Println("input close")
}

func (s *StdInIo) Close() {
	go func() {
		s.IsCloseChan <- true
	}()
}

func NewStdInIo() *StdInIo {
	return &StdInIo{
		IsCloseChan: make(chan bool),
	}
}
