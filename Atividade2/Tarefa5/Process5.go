package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
	"bufio"
	"encoding/json"
)

var err error
var myProcess int
var myPort string
var allPorts [4] string
var nServers int
var CliConn [3]*net.UDPConn
var ServerAddr *net.UDPAddr
var ServConn *net.UDPConn

type VectorClock struct {
	MyProc int
	AllProcessesClocks []int
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

    n,_,err := ServConn.ReadFromUDP(buf[0:])
    if err != nil {
        fmt.Println("Error: ",err)
    }

    var receivedVectorClock VectorClock
    err = json.Unmarshal(buf[:n], &receivedVectorClock)
    CheckError(err)
    //updateClock(&myVectorClock, &receivedVectorClock)
    fmt.Println(receivedVectorClock.MyProc)
	fmt.Printf("Meu novo vector clock: (%d", receivedVectorClock.AllProcessesClocks[1])
	for process := 2; process < nServers; process++ {
		fmt.Printf(", ")
		fmt.Printf("%d", receivedVectorClock.AllProcessesClocks[process])
	}
	fmt.Printf(")\n")
    //myVectorClock = receivedVectorClock
}

func doClientJob(otherProcess int, vecClock VectorClock) {
	fmt.Println(vecClock.MyProc)
	fmt.Println("vecclockallpr", vecClock.AllProcessesClocks[1])
	fmt.Println("vecclockallpr", vecClock.AllProcessesClocks[2])
    jsonRequest, err := json.Marshal(vecClock)
    os.Stdout.Write(jsonRequest)
    _,err = CliConn[otherProcess].Write(jsonRequest)
    if err != nil {
        fmt.Println("Error: ", err)
    }
    time.Sleep(time.Second * 1)
}

func initConnections() {
	myProcess, err = strconv.Atoi(os.Args[1])
	CheckError(err)
	myPort = os.Args[myProcess + 1]
	nServers = len(os.Args) - 2
	for i:=1; i <= nServers; i++ {
		allPorts[i] = os.Args[i + 1]
	}

	ServerAddr, err = net.ResolveUDPAddr("udp", ":" + myPort)
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
	myVectorClock := VectorClock {
		myProcess,
		make([]int, nServers + 1),
	}
	fmt.Println(myVectorClock.MyProc)
	myVectorClock.AllProcessesClocks[myProcess] =1

	for {
		go doServerJob()
		select {
		case x, valid := <-ch:
			if valid {
				if x == strconv.Itoa(myProcess) {
					myVectorClock.AllProcessesClocks[myProcess]++
					fmt.Printf("Meu novo vector clock: (%d", myVectorClock.AllProcessesClocks[1])
					for process := 2; process <= nServers; process++ {
						fmt.Printf(", ")
						fmt.Printf("%d", myVectorClock.AllProcessesClocks[process])
					}
					fmt.Printf(")\n");
				} else if _, err := strconv.Atoi(x); err == nil {
					fmt.Printf("Estou enviando meu vector clock\n")
					x, err := strconv.Atoi(x)
					CheckError(err)
					go doClientJob(x, myVectorClock)
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