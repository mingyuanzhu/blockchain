package main

import (
	"bufio"
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
	"io"
	"log"
	"net"
	"strconv"
	"time"
)

const tipInfo = "Enter a new BPM:"

type TcpServer struct {
	port string
}

var bcServer = make(chan []Block)

func (s *TcpServer) run() error {

	server, err := net.Listen("tcp", ":"+s.port)

	if err != nil {
		log.Fatalf("launch tcp server error, %v", err)
		return err
	}

	defer server.Close()

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Fatalf("receive connetion error, %v", err)
		}

		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {

	defer conn.Close()

	io.WriteString(conn, tipInfo)

	scanner := bufio.NewScanner(conn)

	go handleInput(scanner, conn)

	go simulateReceiveBroadcast(conn)

	for range bcServer {
		spew.Dump(Blockchain)
	}

}

// take in BPM from stdin and add it to blockchain after conducting necessary validation
func handleInput(scanner *bufio.Scanner, conn net.Conn) {
	for scanner.Scan() {
		BPM, err := strconv.Atoi(scanner.Text())
		if err != nil {
			log.Printf("input value %s is not a number", scanner.Text())
			continue
		}

		newBlock, err := generateNewBlock(Blockchain[len(Blockchain)-1], BPM)
		if err != nil {
			log.Printf("create block error, %v", err)
		}

		if isBlockValid(Blockchain[len(Blockchain)-1], newBlock) {
			replaceChain(appendChain(newBlock))
		}

		bcServer <- Blockchain

		io.WriteString(conn, tipInfo)
	}
}

// simulate receiving broadcast
func simulateReceiveBroadcast(conn net.Conn) {
	ticker := time.NewTicker(time.Second * 30)
	for range ticker.C {
		mutex.Lock()
		output, err := json.Marshal(Blockchain)
		if err != nil {
			log.Fatalf("sync fail, %v", err)
		}
		mutex.Unlock()
		io.WriteString(conn, string(output))
	}
}
