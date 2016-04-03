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
		DefaultBufferSize = 1024
		DefaultInterval   = 0
	)

	var (
		bufSize    int
		interval   int
		listenAny  bool
		listenIP   string
		listenPort string
		timestamp  bool
	)

	flag.IntVar(&bufSize, "buffer-size", DefaultBufferSize, "Buffer size")
	flag.IntVar(&bufSize, "b", DefaultBufferSize, "Buffer size")

	flag.IntVar(&interval, "interval", DefaultInterval, "Interval of output (millisecond) (default 0)")
	flag.IntVar(&interval, "i", DefaultInterval, "Shorthand of -interval")

	flag.BoolVar(&listenAny, "listen-any", false, "Listen all avalable IP addresses (i.e. 0.0.0.0)")
	flag.BoolVar(&listenAny, "a", false, "Shorthand of -listen-any")

	flag.StringVar(&listenPort, "listen-port", "", "Listen port number")
	flag.StringVar(&listenPort, "p", "", "Shorthand of -listen-port")

	flag.BoolVar(&timestamp, "show-timestamp", false, "Show timestamp")
	flag.BoolVar(&timestamp, "t", false, "Shorthand of -show-timestamp")

	flag.Parse()

	if listenPort == "" {
		fmt.Fprintln(os.Stderr, "udp-viewer: ERROR: Please specify a port number (e.g. --listen-port=4000)")
		os.Exit(1)
	}

	if listenAny {
		listenIP = "0.0.0.0"
	} else {
		listenIP = "127.0.0.1"
	}

	localEP, err := net.ResolveUDPAddr("udp", listenIP+":"+listenPort)
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
	var t string

	for {
		n, remoteEP, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Fprintln(os.Stderr, "udp-viewer: ERROR:", err)
		}
		if timestamp {
			t = time.Now().Format("2006-01-02 15:04:05.00 ")
		}
		fmt.Printf("%s[%s] %s\n", t, remoteEP, string(buf[0:n]))
		time.Sleep(time.Duration(interval) * time.Millisecond) // FIXME
	}
}
