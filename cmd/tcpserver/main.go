package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync/atomic"

	"github.com/godruoyi/go-snowflake"
)

var (
	listen   string
	asBinary bool

	chSendText  = make(chan []byte, 100)
	chSkip      = make(chan any, 100)
	chTerminate = make(chan any, 100)

	currentClient atomic.Uint64
)

func main() {
	flag.StringVar(&listen, "listen", "127.0.0.1:12345", "listen address of server")
	flag.BoolVar(&asBinary, "binary", false, "data print as binary, not text")

	flag.Parse()

	go func() {
		listener, err := net.Listen("tcp", listen)
		if err != nil {
			log.Fatal(err)
		}

		defer listener.Close()

		for {
			conn, e := listener.Accept()
			if e != nil {
				log.Fatal(e)
			}

			processAConnect(conn, asBinary)
		}
	}()

	loop := true

	reader := bufio.NewReader(os.Stdin)
	for loop {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Println("ReadString failed with ", err)

			continue
		}

		if strings.HasSuffix(line, "~~D") {
			log.Println("cancel current input for ~~D")

			continue
		}

		if !strings.HasPrefix(line, "$") {
			if currentClient.Load() == 0 {
				log.Println("!no client, cancel this operation")
				continue
			}

			if strings.HasSuffix(line, "<-\n") {
				chSendText <- []byte(line[:len(line)-3])
			} else {
				chSendText <- []byte(line)
			}

			continue
		}

		line = line[1:]
		line = strings.TrimRight(line, "\r\b\t ")

		if line == "e" || line == "exit" {
			loop = false

			break
		}

		if line == "s" || line == "skip" {
			if currentClient.Load() == 0 {
				log.Println("!no client, cancel this operation")
				continue
			}

			chSkip <- true

			continue
		}

		if line == "t" || line == "terminate" {
			if currentClient.Load() == 0 {
				log.Println("!no client, cancel this operation")
				continue
			}

			chTerminate <- true

			continue
		}

		log.Println("unknown command:", line)
	}
}

func processAConnect(conn net.Conn, asBinary bool) {
	uuid := snowflake.ID()

	currentClient.Store(uuid)

	fnPrint := func(a ...any) {
		aa := []any{uuid, " => "}
		aa = append(aa, a...)

		fmt.Println(aa...)
	}

	fnPrint("new client")

	var terminating bool

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		defer cancel()

		buf := make([]byte, 40960)
		for {
			n, err := conn.Read(buf[:])
			if terminating {
				break
			}

			if err != nil {
				fnPrint("read error:", err)

				break
			}

			fnPrint("receive bytes:", n)

			if asBinary {
				fnPrint("bytes", buf[:n])
			} else {
				fnPrint("text", string(buf[:n]))
			}
		}
	}()

	loop := true

	for loop {
		select {
		case <-ctx.Done():
			log.Println()
			loop = false
			break
		case t := <-chSendText:
			_, _ = conn.Write(t)
		case <-chSkip:
			fnPrint("skip current")
			loop = false
			break
		case <-chTerminate:
			fnPrint("try terminate current")
			terminating = true
			_ = conn.Close()
			break
		}
	}

	currentClient.Store(0)
}
