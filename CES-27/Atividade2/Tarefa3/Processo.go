package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
	"bufio"
)

var err string
var myPort string
var nServers int
var CliConn []*net.UDPConn

var ServConn *net.UDPConn

func CheckError(err error) {
    if err  != nil {
        fmt.Println("Error: ", err)
    }
}

func readInput(ch chan string) {
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _, _ := reader.ReadLine()
		ch <- string(text)
	}
}

func doServerJob() {
    buf := make([]byte, 1024)

    n,addr,err := ServConn.ReadFromUDP(buf)
    fmt.Println("Received ",string(buf[0:n]), " from ",addr)

    if err != nil {
        fmt.Println("Error: ",err)
    }
}

func doClientJob(otherProcess int, i int) {
    msg := strconv.Itoa(i)
    i++
    buf := []byte(msg)
    _,err := CliConn[otherProcess].Write(buf)
    if err != nil {
        fmt.Println(msg, err)
    }
    time.Sleep(time.Second * 1)
}

func initConnections() {
	myPort = os.Args[1]
	nServers = len(os.Args) - 2

	CliConn = make([]*net.UDPConn, nServers)

	ServerAddr, err := net.ResolveUDPAddr("udp", ":" + myPort)
	CheckError(err)
	ServConn, err = net.ListenUDP("udp", ServerAddr)
	CheckError(err)

	LocalAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
    CheckError(err)
	for i:=0; i < nServers; i++ {
		otherServerAddr,err := net.ResolveUDPAddr("udp","127.0.0.1:" + os.Args[i + 2])
	    CheckError(err)

	    CliConn[i], err = net.DialUDP("udp", LocalAddr, otherServerAddr)
	    CheckError(err)
	}
}

func main() {
	initConnections()
	defer ServConn.Close()
	for i := 0; i < nServers; i++ {
		defer CliConn[i].Close()
	}

	i := 0
	ch := make(chan string)
	go readInput(ch)

	for {
		go doServerJob()
		select {
		case x, valid := <-ch:
			if valid {
				fmt.Printf("Recebi do teclado: %s \n", x)
			} else {
				fmt.Println("Channel closed!")
			}
		default:
			time.Sleep(time.Second * 1)
		}
		for j := 0; j < nServers; j++ {
			go doClientJob(j, i)
		}
		time.Sleep(time.Second * 1)
		i++
	}
}