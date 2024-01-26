package main

import (
	"fmt"
	"net"
)

func main() {
	listener, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 6999, //對外公開 Prt 5000
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Local: <%s> \n", listener.LocalAddr().String())

	//定義有兩台電腦要互相連接
	var a_computer, b_computer net.UDPAddr

	//等待第一台電腦連接
	a_computer = *waitAndRead(listener)
	//等待第二台電腦連接
	b_computer = *waitAndRead(listener)

	//把 b 電腦資訊傳給 a 電腦
	listener.WriteToUDP([]byte(b_computer.String()), &a_computer)
	//把 a 電腦資訊傳給 b 電腦
	listener.WriteToUDP([]byte(a_computer.String()), &b_computer)

	//伺服器可以掰掰了
	fmt.Println("Server Mission Clear.")
}

func waitAndRead(listener *net.UDPConn) *net.UDPAddr {

	buff := make([]byte, 1024)
	n, remoteAddr, err := listener.ReadFromUDP(buff)
	if err != nil {
		fmt.Printf("error during read %s", err)
	}

	fmt.Printf("<%s> %s\n", remoteAddr, buff[:n])

	//port := strconv.Itoa(remoteAddr.Port)
	return remoteAddr

}
