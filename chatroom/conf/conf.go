package conf

import "fmt"

const (
	UDPSERVERIP = "0.0.0.0"
	UDPSERVERPORT = 30310

	TCP_SERVER_IP = "0.0.0.0"
	TCP_SERVER_PORT = 30311
)

func GetUdpServerAddress() string {
	 return fmt.Sprintf("%s:%d",UDPSERVERIP,UDPSERVERPORT)
}

func GetTcpServerAddress() string  {
    return fmt.Sprintf("%s:%d",TCP_SERVER_IP,TCP_SERVER_PORT)
}

