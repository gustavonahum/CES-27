package main


import (
	"fmt"
	"net"
	"os"
	"strings"
	"strconv"
	"time"
	"bufio"
)


type Process struct {
	Id int
	Participant bool
}


var err error
var nProcesses int
var processes []Process
var nElectionBeginners int
var electionBeginners [] int

var ServerAddr []*net.UDPAddr
var ServerConn []*net.UDPConn
var CliAddr []*net.UDPAddr
var CliConn []*net.UDPConn


func CheckError(err error) {
    if err  != nil {
        fmt.Println("Error: ", err)
    }
}


func next(pIndex int) int {
	return (pIndex + 1) % nProcesses
}


func initialize() {
	// Quantidade de processos
	nProcesses = len(os.Args) - 1

	// Cria os processos
	processes = make([]Process, nProcesses)
	for i:=0; i < nProcesses; i++ {
	    processes[i].Id, err = strconv.Atoi(os.Args[i + 1])
	    CheckError(err)
		processes[i].Participant = false
	}

	// Monta as conexoes UDP
	ServerAddr = make([]*net.UDPAddr, nProcesses)
	ServerConn = make([]*net.UDPConn, nProcesses)

	CliAddr = make([]*net.UDPAddr, nProcesses)
	CliConn = make([]*net.UDPConn, nProcesses)

	for i:=0; i < nProcesses; i++ {
		ServerAddr[i], err = net.ResolveUDPAddr("udp","127.0.0.1:" + os.Args[i + 1])
	    CheckError(err)
		ServerConn[i], err = net.ListenUDP("udp", ServerAddr[i])
		CheckError(err)
	}

	for i:=0; i<nProcesses; i++ {
		CliAddr[i], err = net.ResolveUDPAddr("udp", "127.0.0.1:0")
		CheckError(err)
	    CliConn[i], err = net.DialUDP("udp", CliAddr[i], ServerAddr[next(i)])
	    CheckError(err)
	}
}


func readInput(ch chan string) {
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _, _ := reader.ReadLine()
		ch <- string(text)
	}
}


func sendMessage(pIndex int, message string) {
	fmt.Printf("P%d: Enviou %s \n", processes[pIndex].Id, message)
	if strings.Index(message, "S") == 0 {
		processes[pIndex].Participant = true
		time.Sleep(time.Second * 1)
		go doClientJob(pIndex, message)
	} else if strings.Index(message, "F") == 0 {
		processes[pIndex].Participant = false
		time.Sleep(time.Second * 1)
		go doClientJob(pIndex, message)
	} else {
		fmt.Printf("Falha na comunicacao \n")
	}
}


func beginElection(pId int) {
	for pIndex:=0; pIndex < nProcesses; pIndex++ {
		if processes[pIndex].Id == pId {
			sendMessage(pIndex, "S" + strconv.Itoa(pId))
			return
		}
	}
	fmt.Printf("O numero de processo inserido esta incorreto! \n")
}


func doServerJob(pIndex int) {
    buf := make([]byte, 1024)

    n,_,err := ServerConn[pIndex].ReadFromUDP(buf)
    if err != nil {
        fmt.Println("Error: ",err)
    }

    fmt.Printf("P%d: Recebeu %s \n",processes[pIndex].Id,string(buf[0:n]))

    if (string(buf[0]) == "S") {
    	strPId := string(buf[1:n])
	    intPId, err := strconv.Atoi(strPId)
	    CheckError(err)
	    if (intPId > processes[pIndex].Id) {
	    	// Marca-se como participante, e encaminha a mensagem
	    	processes[pIndex].Participant = true
	    	sendMessage(pIndex, string(buf[0:n]))
	    } else if (intPId < processes[pIndex].Id) {
	    	// Descarta a mensagem, e inicia a eleicao, caso ainda nao o tenha feito
	    	if (!processes[pIndex].Participant) {
		    	processes[pIndex].Participant = true
		    	sendMessage(pIndex, "S" + strconv.Itoa(processes[pIndex].Id))
	    	}
	    } else {
	    	// Processo pIndex torna-se o lider
	    	processes[pIndex].Participant = false
	    	sendMessage(pIndex, "F" + strconv.Itoa(processes[pIndex].Id))
	    }
    } else if (string(buf[0]) == "F") {
    	strPId := string(buf[1:n])
	    intPId, err := strconv.Atoi(strPId)
	    CheckError(err)
	    if (intPId == processes[pIndex].Id) {
	    	fmt.Printf("O processo %d foi eleito como lider! \n", processes[pIndex].Id)
	    } else {
			processes[pIndex].Participant = false
		    sendMessage(pIndex, string(buf[0:n]))
	    }
    }
}


func doClientJob(pIndex int, message string) {
    buf := []byte(message)
    _,err := CliConn[pIndex].Write(buf)
    if err != nil {
        fmt.Println(message, err)
    }
}


func main() {
	initialize()

	for i:=0; i < nProcesses; i++ {
		defer ServerConn[i].Close()
		defer CliConn[i].Close()
	}

	ch := make(chan string, 100)
	go readInput(ch)

	for {
		for i:=0; i < nProcesses; i++ {
			go doServerJob(i)
		}

		select {
		case x, valid := <-ch:
			if valid {
				if strings.Index(x, "start") == 0 {
					inputSlice := strings.Split(x, " ")
					for i:=1; i < len(inputSlice); i++ {
						pId, err := strconv.Atoi(inputSlice[i])
	    				CheckError(err)
	    				go beginElection(pId)
					}
				} else {
					fmt.Printf("O formato do input esta incorreto! \n")
				}
			} else {
				fmt.Println("Channel closed!")
			}
		default:
			time.Sleep(time.Second * 1)
		}
	}
}