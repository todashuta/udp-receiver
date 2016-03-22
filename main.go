package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	const (
		DefaultAddress    = "127.0.0.1"
		DefaultBufferSize = 1024
		DefaultPort       = "8080"
		DefaultInterval   = 0
	)

	var (
		addr     string
		bufSize  int
		port     string
		interval int
	)

	flag.StringVar(&addr, "addr", DefaultAddress, "Server address")
	flag.StringVar(&addr, "a", DefaultAddress, "Shorthand of -addr")
	flag.IntVar(&bufSize, "buffer-size", DefaultBufferSize, "Buffer size")
	flag.IntVar(&bufSize, "b", DefaultBufferSize, "Buffer size")
	flag.StringVar(&port, "port", DefaultPort, "Port number")
	flag.StringVar(&port, "p", DefaultPort, "Shorthand of -port")
	flag.IntVar(&interval, "interval", DefaultInterval, "Interval of output (millisecond) (default 0)")
	flag.IntVar(&interval, "i", DefaultInterval, "Shorthand of -interval")
	flag.Parse()

	serverAddr, err := net.ResolveUDPAddr("udp", addr+":"+port)
	if err != nil {
		fmt.Fprintln(os.Stderr, "udp-viewer: ERROR:", err)
		os.Exit(1)
	}

	conn, err := net.ListenUDP("udp", serverAddr)
	if err != nil {
		fmt.Fprintln(os.Stderr, "udp-viewer: ERROR:", err)
		os.Exit(1)
	}
	defer conn.Close()

	buf := make([]byte, bufSize)

	for {
		n, remoteAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Fprintln(os.Stderr, "udp-viewer: ERROR:", err)
		}
		fmt.Printf("[%s] %s\n", remoteAddr, string(buf[0:n]))
		time.Sleep(time.Duration(interval) * time.Millisecond)
	}
}
