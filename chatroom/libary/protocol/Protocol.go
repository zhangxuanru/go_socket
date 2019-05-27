package protocol

import (
	"bytes"
	"encoding/binary"
)

const (
	HEADERSTR = "Zhangxr"
	HENDERLEN = 7
	BINARYLEN = 4
)

func Pack(data []byte) []byte  {
	  length := len(data)
      message := append([]byte(HEADERSTR),IntToBytes(length)...)
      return append(message,data...)
}


func Unpack(data []byte,msgChan chan string) []byte  {
	  var(
	     i int
	     length int
		 messageLen int
	  )
      length = len(data)
      if length < HENDERLEN+BINARYLEN{
       	  return data[0:]
	  }
      for i=0; i<length;i++{
      	 if length < i+HENDERLEN+BINARYLEN{
      	 	break
		 }
		 if string(data[i:i+HENDERLEN]) == HEADERSTR{
              messageLen = BytesToInt(data[i+HENDERLEN:i+HENDERLEN+BINARYLEN])
              if length < i+HENDERLEN+BINARYLEN+messageLen{
                	break
			  }
			  message := data[i+HENDERLEN+BINARYLEN:i+HENDERLEN+BINARYLEN+messageLen]
			  msgChan <- string(message)
			  i = i+HENDERLEN+BINARYLEN+messageLen-1
		 }
	  }
	if i == length {
		return make([]byte, 0)
	}
    return data[i:]
}

func IntToBytes(n int) []byte  {
    x := int32(n)
    buf := bytes.NewBuffer([]byte{})
    binary.Write(buf,binary.BigEndian,x)
    return buf.Bytes()
}


func BytesToInt(buf []byte) int {
	 var x int32
	 byteBuffer := bytes.NewBuffer(buf)
     binary.Read(byteBuffer,binary.BigEndian,&x)
	 return int(x)
}

