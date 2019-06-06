package common

import (
	"fmt"
	"socket/app/config"
)

func GetServerAddress() string {
	return fmt.Sprintf("%s:%d", config.NETWORK_SERVER_IP, config.NETWORK_SERVER_PORT)
}

func GetNetWorkType() string {
	return config.NETWORK_TYPE
}
