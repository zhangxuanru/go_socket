package main

import  "socket/chatroom/libary/chat/client"

func main()  {
	client.InitClient()

	//var(
	//	udpAddr *net.UDPAddr
	//	udpConn *net.UDPConn
	//	err error
	//)
	//if udpAddr,err = net.ResolveUDPAddr("udp",conf.GetUdpServerAddress());err!=nil{
	//	fmt.Println("error:",err.Error())
	//	return
	//}
	//if udpConn,err = net.DialUDP("udp",nil,udpAddr);err!=nil{
	//	fmt.Println("error:",err.Error())
	//	return
	//}
	//for{
	//	buf := make([]byte,1024)
	//	fmt.Println("input msg:")
	//	msg,_:= bufio.NewReader(os.Stdin).ReadString('\n')
	//	udpConn.Write([]byte(msg))
	//
	//	udpConn.Read(buf)
	//	fmt.Println("read server msg:",string(buf))
	//}

}
