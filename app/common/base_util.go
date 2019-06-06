package common

import (
	"fmt"
	"strings"
)

func PrintErr(err error) {
	fmt.Println("err:", err)
}

func IsMsgBye(msg string) bool {
	return strings.EqualFold(msg, "bye")
}
