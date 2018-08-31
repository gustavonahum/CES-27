package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
	"bufio"
)

var err error
var myProcess int
var myPort string
var allPorts [] string
var nServers int
var CliConn []*net.UDPConn
var ServConn *net.UDPConn

var logicalClock int

func updateClock(clock1 int, clock2 int) int {
	if clock1 > clock2 {
		return clock1 + 1
	}
	return clock2 + 1
}

func CheckError(err error) {
    if err  != nil {
        fmt.Println("Error: ", err)
    }
}

func PrintError(err error) {
	if err != nil {
		fmt.Println("Erro: ", err)
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

    n,_,err := ServConn.ReadFromUDP(buf)
    if err != nil {
        fmt.Println("Error: ",err)
    }
    fmt.Printf("Estou recebendo o seguinte logical clock: %s \n", string(buf[0:n]))
    receivedClock, err := strconv.Atoi(string(buf[0:n]))
    CheckError(err)
    logicalClock = updateClock(logicalClock, receivedClock)
    fmt.Printf("Meu novo logical clock eh: %d \n", logicalClock)
}

func doClientJob(otherProcess int, logClock int) {
    msg := strconv.Itoa(logClock)
    buf := []byte(msg)
    _,err := CliConn[otherProcess].Write(buf)
    if err != nil {
        fmt.Println(msg, err)
    }
    time.Sleep(time.Second * 1)
}

func initConnections() {
	myProcess, err = strconv.Atoi(os.Args[1])
	CheckError(err)
	myPort = os.Args[myProcess + 1]
	nServers = len(os.Args) - 2

	allPorts = make([]string, nServers + 1)
	for i:=1; i <= nServers; i++ {
		allPorts[i] = os.Args[i + 1]
	}

	CliConn = make([]*net.UDPConn, nServers + 1)

	ServerAddr, err := net.ResolveUDPAddr("udp", ":" + myPort)
	CheckError(err)
	ServConn, err = net.ListenUDP("udp", ServerAddr)
	CheckError(err)

	LocalAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
    CheckError(err)
	for i:=1; i <= nServers; i++ {
		if i != myProcess {
			otherServerAddr,err := net.ResolveUDPAddr("udp","127.0.0.1:" + allPorts[i])
		    CheckError(err)

		    CliConn[i], err = net.DialUDP("udp", LocalAddr, otherServerAddr)
		    CheckError(err)
		}
	}
}

func main() {
	initConnections()
	defer ServConn.Close()
	for i := 1; i <= nServers; i++ {
		if i != myProcess {
			defer CliConn[i].Close()
		}	
	}

	i := 0
	ch := make(chan string)
	go readInput(ch)
	logicalClock = 1

	for {
		go doServerJob()
		select {
		case x, valid := <-ch:
			if valid {
				if x == strconv.Itoa(myProcess) {
					logicalClock++
					fmt.Printf("Meu novo logical clock: %d \n", logicalClock)
				} else if _, err := strconv.Atoi(x); err == nil {
					fmt.Printf("Estou enviando meu logical clock: %d \n", logicalClock)
					x, err := strconv.Atoi(x)
					CheckError(err)
					go doClientJob(x, logicalClock)
				}
			} else {
				fmt.Println("Channel closed!")
			}
		default:
			time.Sleep(time.Second * 1)
		}
		time.Sleep(time.Second * 1)
		i++
	}
}