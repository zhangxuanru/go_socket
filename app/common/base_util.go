package common

import (
	"fmt"
	"strings"
	"bytes"
)

func PrintErr(err error) {
	fmt.Println("err:", err)
}

func IsMsgBye(msg string) bool {
	return strings.EqualFold(msg, "bye")
}

func ByteToString(buf []byte) string {
	index := bytes.IndexByte(buf,0)
	if index > -1{
	    return strings.TrimSpace(string(buf[:index]))
	}
	 return strings.TrimSpace(string(buf))
 }
