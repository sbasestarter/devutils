package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

// -listen :1112 -remote remote-machine:22

func main() {
	var listen string
	var remoteAddr string

	flag.StringVar(&listen, "listen", ":10001", "listen address")
	flag.StringVar(&remoteAddr, "remote", "", "remote address")
	flag.Parse()

	if listen == "" {
		log.Fatal("no listen address")
	}

	if remoteAddr == "" {
		log.Fatal("no remote address")
	}

	tlAddr, err := net.ResolveTCPAddr("tcp4", listen)
	if err != nil {
		log.Fatal("resolve tcp addr failed:", err)
	}

	trAddr, err := net.ResolveTCPAddr("tcp4", remoteAddr)
	if err != nil {
		log.Fatal("resolve remote tcp addr failed:", err)
	}

	listener, err := net.ListenTCP("tcp", tlAddr)
	if err != nil {
		log.Fatal("listen tcp failed:", err)
	}

	fmt.Printf("server listen on %s\n", listen)

	for {
		conn, e := listener.AcceptTCP()
		if e != nil {
			log.Println("accept failed:", e)

			continue
		}

		go dealClientConn(conn, trAddr)
	}
}

func dealClientConn(conn *net.TCPConn, remoteAddr *net.TCPAddr) {
	defer conn.Close()

	info := fmt.Sprintf("%s <-> %s", conn.RemoteAddr().String(), remoteAddr.String())

	remoteConn, err := net.DialTCP("tcp", nil, remoteAddr)
	if err != nil {
		log.Printf("%s: dial failed: %v\n", info, err)

		return
	}

	err = conn.SetKeepAlive(true)
	if err != nil {
		log.Printf("%s: set local tcp keep-alive failed: %s\n", info, err.Error())
	}

	err = conn.SetKeepAlivePeriod(time.Second * 10) //nolint:mnd // .
	if err != nil {
		log.Printf("%s: set local tcp keep-alive-period failed: %s\n", info, err.Error())
	}

	err = remoteConn.SetKeepAlive(true)
	if err != nil {
		log.Printf("%s: set remote tcp keep-alive failed: %s\n", info, err.Error())
	}

	err = remoteConn.SetKeepAlivePeriod(time.Second * 10) //nolint:mnd // .
	if err != nil {
		log.Printf("%s: set remote tcp keep-alive-period failed: %s\n", info, err.Error())
	}

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer func() {
			_ = remoteConn.Close()
			_ = conn.Close()

			wg.Done()
		}()

		_, e := io.CopyBuffer(remoteConn, conn, nil)

		log.Printf("%s: x data local to remote end: %v\n", info, e)
	}()

	_, _ = io.CopyBuffer(conn, remoteConn, nil)

	log.Printf("%s: x data remote to local end: %v\n", info, err)

	_ = remoteConn.Close()
	_ = conn.Close()

	wg.Wait()

	log.Printf("%s: deal client connection finished\n", info)
}
