package util

import (
	"bytes"
	"strings"
)

func ByteToString(buf []byte) string {
	index := bytes.IndexByte(buf,0)
	return strings.TrimSpace(string(buf[:index]))
}

