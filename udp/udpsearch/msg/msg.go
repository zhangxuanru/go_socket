package msg

import (
	"fmt"
	"strings"
	"errors"
	"strconv"
)

var(
	SN_HEADER = "收到暗号，我是(SN):"
	PORT_HEADER = "这是暗号,请回电商品(PORT):"
)

func BuildWithPort(port int ) string  {
    return fmt.Sprintf(PORT_HEADER+"%d",port)
}

func ParsePort(str string) (int,error)  {
	var(
		index int
		strPort string
	)
	  if index = strings.LastIndex(str,PORT_HEADER); index > -1 {
		   strPort = strings.TrimPrefix(str,PORT_HEADER)
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

