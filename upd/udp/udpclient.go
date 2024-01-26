package main

import (
	"fmt"
	"net"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func main() {

	//選擇一個本地打洞的 port，我在這裡是固定為 5001 (注意，那表示下面連接時是固定連到對方的 5001，你不能在自己電腦開兩個 5001 的 app)
	srcAddr := &net.UDPAddr{IP: net.IPv4zero, Port: 5001} //fixed same port for all computer, do not use 0 by dynamic

	//輸入公開網路那台 server 的 ip 以及那個 port (請替換 0.0.0.0)
	//dstAddr := &net.UDPAddr{IP: net.ParseIP("0.0.0.0"), Port: 5000}
	dstAddr := &net.UDPAddr{IP: net.ParseIP("8.218.116.95"), Port: 6999}

	//先跟 server 連接
	conn, err := net.DialUDP("udp", srcAddr, dstAddr)
	if err != nil {
		fmt.Println(err)
	}

	//隨便傳一個訊息給 server ，目的是讓 server 知道我的 IP 地址
	//send msg to central public server(let server know my address)
	conn.Write([]byte("hello"))
	fmt.Println("Hole Punching, UDP Client to Public Server.")

	//強制等待 server 告訴我對方的 IP
	//receive another computer information...
	fmt.Println("Waiting for Public Server Return Message...")
	var another_computer *net.UDPAddr = waitAndReadAndParseUDP(conn)

	fmt.Println("Successfuly get another computer address!")
	fmt.Println("ADDRESS: ", another_computer.String())

	//連接成功，跟 server 斷線，準備跟另一台電腦連接
	fmt.Println("Close Original UDP Connection(Port release)")
	conn.Close()

	//直接連接對方電腦，使用的 port 是 5001 跟自己一樣
	//send message to other computer directly
	fmt.Println("Build a connection to another computer port and OPEN NAT PORT wait for receive message")
	another_conn, err := net.DialUDP("udp", srcAddr, another_computer)
	if err != nil {
		fmt.Println("Can't Build Hole Punching...")
		fmt.Println(err)
	}
	defer another_conn.Close()
	_, err = conn.Write([]byte("Hole Punching"))
	if err != nil {
		fmt.Println("Failed to send punched message")
		fmt.Println(err)
	}

	//建立連接之後，瘋狂送訊息
	fmt.Println("Keep Send Message")
	//keep send message...
	var count int = 0
	go func() {
		for {
			msg := fmt.Sprintf("%s  message %d\n", runtime.GOOS, count)
			c, err := another_conn.Write([]byte(msg))
			if err != nil {
				fmt.Println("Failed to send message")
				fmt.Println(err)
			} else {
				fmt.Print(msg + "  write " + strconv.Itoa(int(c)))
			}
			time.Sleep(5 * time.Second)

			count++
		}
	}()

	//瘋狂接收對方訊息
	//keep receiving message
	for {
		buff := make([]byte, 1024)
		n, remoteAddr, err := another_conn.ReadFromUDP(buff)
		if err != nil {
			fmt.Printf("error during read %s", err)
		}

		fmt.Printf("receive from <%s> %s\n", remoteAddr, buff[:n])
	}
}

func waitAndReadAndParseUDP(listener *net.UDPConn) *net.UDPAddr {

	buff := make([]byte, 1024)
	n, _, err := listener.ReadFromUDP(buff)
	if err != nil {
		fmt.Printf("error during read %s", err)
	}

	another_computer_address := strings.Split(string(buff[:n]), ":")
	port, _ := strconv.Atoi(another_computer_address[1])

	return &net.UDPAddr{
		IP:   net.ParseIP(another_computer_address[0]),
		Port: port,
	}

}
