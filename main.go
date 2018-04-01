package main

import (
	"flag"
	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var httpMode = flag.Bool(
	"httpMode",
	true,
	"Start http server, if the value is false will use tcp mode",
)

var tcpMode = flag.Bool(
	"tcpMode",
	true,
	"Start tcp server, if the value is false will use tcp mode",
)

const (
	httpAddr = "HTTP_ADDR"
	tcpAddr  = "TCP_ADDR"
)

var mutex = &sync.Mutex{}

func main() {

	flag.Parse()

	err := godotenv.Load()

	if err != nil {
		log.Fatal(err)
		return
	}

	go genesisBlock()

	if *httpMode {
		s := &HttpServer{
			port: os.Getenv(httpAddr),
		}
		go func() {
			log.Fatal(s.run())
		}()
	}

	if *tcpMode {
		s := &TcpServer{
			port: os.Getenv(tcpAddr),
		}
		go func() {
			log.Fatal(s.run())
		}()
	}

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	log.Printf("Terminate by signal:%v", <-ch)
}
func genesisBlock() {
	t := time.Now()
	genesisBlock := Block{0, t.String(), 0, "", ""}
	spew.Dump(genesisBlock)
	Blockchain = appendChain(genesisBlock)
}
