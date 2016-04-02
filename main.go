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
		DefaultHost       = "127.0.0.1"
		DefaultBufferSize = 1024
		DefaultPort       = "8080"
		DefaultInterval   = 0
	)

	var (
		host     string
		bufSize  int
		port     string
		interval int
	)

	flag.StringVar(&host, "host", DefaultHost, "Sender host address")
	flag.StringVar(&host, "h", DefaultHost, "Shorthand of -host")
	flag.IntVar(&bufSize, "buffer-size", DefaultBufferSize, "Buffer size")
	flag.IntVar(&bufSize, "b", DefaultBufferSize, "Buffer size")
	flag.StringVar(&port, "port", DefaultPort, "Sender port number")
	flag.StringVar(&port, "p", DefaultPort, "Shorthand of -port")
	flag.IntVar(&interval, "interval", DefaultInterval, "Interval of output (millisecond) (default 0)")
	flag.IntVar(&interval, "i", DefaultInterval, "Shorthand of -interval")
	flag.Parse()

	localEP, err := net.ResolveUDPAddr("udp", host+":"+port)
	if err != nil {
		fmt.Fprintln(os.Stderr, "udp-viewer: ERROR:", err)
		os.Exit(1)
	}

	conn, err := net.ListenUDP("udp", localEP)
	if err != nil {
		fmt.Fprintln(os.Stderr, "udp-viewer: ERROR:", err)
		os.Exit(1)
	}
	defer conn.Close()

	buf := make([]byte, bufSize)

	for {
		n, remoteEP, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Fprintln(os.Stderr, "udp-viewer: ERROR:", err)
		}
		fmt.Printf("[%s] %s\n", remoteEP, string(buf[0:n]))
		time.Sleep(time.Duration(interval) * time.Millisecond)
	}
}
