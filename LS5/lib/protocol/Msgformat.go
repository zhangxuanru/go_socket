package protocol

import (
	"strconv"
	"bytes"
	"fmt"
	"strings"
	"errors"
	"socket/LS5/conf"
)

var(
	SN_HEADER = "收到暗号，我是(SN):"
	PORT_HEADER = "这是暗号,请回电商品(PORT):"
    UDP_MSG_HEADER  = []byte("7,7,7,7,7,7,7")
)

func BuildWithPort(port int ) string  {
	  return fmt.Sprintf(PORT_HEADER+"%d",port)
}

func ParsePort(str string) (int,error)  {
	var(
		index int
		strPort string
	)
	str = strings.TrimSpace(str)
	if index = strings.LastIndex(str,PORT_HEADER); index > -1 {
		pos := index+len(PORT_HEADER)
		strPort = str[pos:]
		return strconv.Atoi(strPort)
	}
	return -1,errors.New("not search PORT_HEADER")
}


func BuildWithSn(sn string) string {
	return  sn+"\n"
}


func ParseSn(str string) string {
	var(
		index int
	)
	if index = strings.LastIndex(str,SN_HEADER); index > -1 {
		return  strings.TrimPrefix(str,SN_HEADER)
	}
	return ""
}


func GetStringByBuF2(buf []byte) string {
	index := bytes.IndexByte(buf,0)
	return strings.TrimSpace(string(buf[:index]))
}


func GetUdpServerAddress() string {
	return fmt.Sprintf("%s:%d",conf.UDP_SERVER_IP,conf.UDP_SERVER_PORT)
}

func GetTcpServerAddress() string {
	return fmt.Sprintf("%s:%d",conf.TCP_SERVER_IP,conf.TCP_SERVER_PORT)
}

func GetUdpBroadcastAddress() string {
	return fmt.Sprintf("%s:%d",conf.SEARCH_UDP_IPADDR,conf.UDP_SERVER_PORT)
}


func CheckReceiveMsgHeader(buf []byte) bool  {
	   headerLen := len(UDP_MSG_HEADER)
	   return bytes.Equal(buf[:headerLen],UDP_MSG_HEADER)
}
