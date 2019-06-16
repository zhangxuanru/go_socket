package trans

import (
	"bytes"
	"encoding/binary"
)

const PACK_DATE_HEADER = "zxrDjGoodGoodGood"
const PACK_BINARY_LEN = 4

//打包
func Pack(data []byte) []byte {
	length := len(data)
	lenBytes := IntToBytes(length)
	message := append([]byte(PACK_DATE_HEADER), lenBytes...)
	return append(message, data...)
}

//解包
func UnPack(data []byte, readChan chan []byte) []byte {
	var i int
	headerLen := len(PACK_DATE_HEADER)
	dataLen := len(data)
	if dataLen < headerLen+PACK_BINARY_LEN {
		return data[0:]
	}
	for i = 0; i < dataLen; i++ {
		if dataLen < i+headerLen+PACK_BINARY_LEN {
			break
		}
		if string(data[i:i+headerLen]) == PACK_DATE_HEADER {
			messageLen := BytesToInt(data[i+headerLen : i+headerLen+PACK_BINARY_LEN])
			if dataLen < i+headerLen+PACK_BINARY_LEN+messageLen {
				break
			}
			message := data[i+headerLen+PACK_BINARY_LEN : i+headerLen+PACK_BINARY_LEN+messageLen]
			readChan <- message
			i = i + headerLen + PACK_BINARY_LEN + messageLen - 1 //由于for循环里有自增， 所以这里先-1，
		}
	}
	if i == dataLen {
		return make([]byte, 0)
	}
	return data[i:]
}

func IntToBytes(n int) []byte {
	x := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, &x)
	return bytesBuffer.Bytes()
}

func BytesToInt(data []byte) int {
	var x int32
	byteBuffer := bytes.NewBuffer(data)
	binary.Read(byteBuffer, binary.BigEndian, &x)
	return int(x)
}
