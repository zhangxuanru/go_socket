package common

import (
	"bytes"
	"fmt"
	"socket/app/config"
	"strings"
)

func PrintErr(err error) {
	fmt.Println("err:", err)
}

func IsMsgBye(msg string) bool {
	return strings.EqualFold(msg, "bye")
}

func IsBufBye(buf []byte) bool {
	msgStr := ByteToString(buf)
	return IsMsgBye(msgStr)
}

func ByteToString(buf []byte) string {
	index := bytes.IndexByte(buf, 0)
	if index > -1 {
		return strings.TrimSpace(string(buf[:index]))
	}
	return strings.TrimSpace(string(buf))
}

func CheckBye(msg []byte, IsCloseChan chan bool) {
	msgStr := ByteToString(msg)
	if IsMsgBye(msgStr) {
		go func() {
			IsCloseChan <- true
		}()
	}
}

func RemoveStrSendHeader(byt []byte) []byte {
	byt = bytes.TrimPrefix(byt, []byte(config.SEND_FILE_HEADER_PACK))
	byt = bytes.TrimPrefix(byt, []byte(config.SEND_STR_HEADER_PACK))
	return byt
}
