package trans

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"os"
	"path"
	"path/filepath"
	"socket/app/common"
	"socket/app/config"
	"strings"
	"time"
)

type TcpFilePack struct {
	closeChan chan bool
	TempFile
}

type TempFile struct {
	tempFile     *os.File
	tempFileName string
}

const TEMPFILEDIR = "/home/work/baihe/go-mooc/temp/"

func (s *TcpFilePack) Write(conn *net.TCPConn, data []byte) (int, error) {
	if conn == nil {
		return 0, errors.New("conn is nil")
	}
	header := []byte(config.SEND_FILE_HEADER_PACK)
	conn.Write(Pack(header))
	if bytes.HasPrefix(data, []byte("file:")) {
		SendFile(conn, data)
	} else {
		conn.Write(Pack(data))
	}
	return 0, nil
}

func (s *TcpFilePack) Read(conn *net.TCPConn, readMsgChan chan []byte) (err error) {
	//NewRead().Read(conn, readMsgChan)
	return nil
}

func (s *TcpFilePack) Close(conn *net.TCPConn) error {
	fmt.Println("conn is close")
	return conn.Close()
}

func (s *TcpFilePack) Receive(receiveMsg []byte, receiveChan chan []byte, source string) {
	fmt.Println("file  receiveMsg:", string(receiveMsg))
	if len(receiveMsg) == 0 || strings.EqualFold(source, config.CLIENTIDENT) {
		return
	}
	//生成临时文件
	if _, err := s.createTempFile(); err != nil {
		fmt.Println("err:", err)
		return
	}
	if bytes.HasPrefix(receiveMsg, []byte("#end:#")) {
		//重命名临时文件
		fileName := bytes.TrimPrefix(receiveMsg, []byte("#end:#"))
		ext := path.Ext(string(fileName))
		fName := path.Base(s.tempFileName)
		fName = strings.Replace(fName, ".temp", ext, -1)
		if err := s.RenameFile(fName); err != nil {
			os.Remove(s.tempFileName)
			go func() {
				receiveChan <- []byte(err.Error())
			}()
		} else {
			go func() {
				receiveChan <- []byte("file save success path:" + s.tempFileName)
				s.tempFile = nil
				s.tempFileName = ""
			}()
		}
	} else {
		//写入临时文件
		if err := s.WriteTempFile(receiveMsg); err != nil {
			os.Remove(s.tempFileName)
			go func() {
				receiveChan <- []byte(err.Error())
			}()
			return
		}
	}
}

func (s *TcpFilePack) createTempFile() (*os.File, error) {
	if s.tempFile != nil {
		return s.tempFile, nil
	}
	tempFileName := time.Now().Format("2006-01-02")
	rand1 := rand.New(rand.NewSource(time.Now().UnixNano()))
	tempFileName = fmt.Sprintf("%s%s%d%s", TEMPFILEDIR, tempFileName, rand1.Int31(), ".temp")
	file, err := os.OpenFile(tempFileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0666)
	s.tempFile = file
	s.tempFileName = tempFileName
	return s.tempFile, err
}

func (s *TcpFilePack) WriteTempFile(data []byte) error {
	writer := bufio.NewWriter(s.tempFile)
	if _, err := writer.Write(data); err != nil {
		return err
	}
	if err := writer.Flush(); err != nil {
		return err
	}
	return nil
}

func (s *TcpFilePack) RenameFile(newName string) error {
	s.tempFile.Close()
	if err := os.Rename(s.tempFileName, TEMPFILEDIR+newName); err != nil {
		return err
	}
	s.tempFileName = TEMPFILEDIR + newName
	return nil
}

func SendFile(conn *net.TCPConn, fileMsg []byte) (err error) {
	var (
		fileInfo *os.File
	)
	fileMsg = common.RemoveFilePrefixMsg(fileMsg)
	filePath := string(fileMsg)
	if common.IsFileNotExists(filePath) == true {
		err = errors.New("file is not exists")
		return
	}
	if fileInfo, err = os.Open(filePath); err != nil {
		return
	}
	_, file := filepath.Split(filePath)
	defer fileInfo.Close()
	buf := make([]byte, 1024)
	for {
		n, e := fileInfo.Read(buf)
		if e == io.EOF {
			break
		}
		conn.Write(Pack(buf[:n]))
	}
	endMsg := []byte("#end:#" + file)
	conn.Write(Pack(endMsg))
	return
}

func NewTcpFilePack() *TcpFilePack {
	return &TcpFilePack{
		closeChan: make(chan bool),
	}
}
