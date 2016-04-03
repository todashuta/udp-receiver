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

	flag.IntVar(&bufSize, "buffer-size", DefaultBufferSize, "")
	flag.IntVar(&bufSize, "b", DefaultBufferSize, "")

	flag.IntVar(&interval, "interval", DefaultInterval, "")
	flag.IntVar(&interval, "i", DefaultInterval, "")

	flag.BoolVar(&listenAny, "listen-any", false, "")
	flag.BoolVar(&listenAny, "a", false, "")

	flag.StringVar(&listenPort, "listen-port", "", "")
	flag.StringVar(&listenPort, "p", "", "")

	flag.BoolVar(&timestamp, "show-timestamp", false, "")
	flag.BoolVar(&timestamp, "t", false, "")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: udp-receiver [OPTION]...
Options:
  -a, -listen-any      : Listen all available IP addresses (i.e. 0.0.0.0)
  -b, -buffer-size     : Buffer size
  -p, -listen-port     : Listen port number
  -t, -show-timestamp  : Show timestamp

Experimental options:
  -i, -interval        : Output interval (millisecond) (default 0)
`)
		os.Exit(2)
	}

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
