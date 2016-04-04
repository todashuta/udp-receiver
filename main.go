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
		bufSize       int
		interval      int
		listenAny     bool
		listenIP      string
		listenPort    string
		showTimestamp bool
		showSender    bool
	)

	flag.IntVar(&bufSize, "buffer-size", DefaultBufferSize, "")
	flag.IntVar(&bufSize, "b", DefaultBufferSize, "")

	flag.IntVar(&interval, "interval", DefaultInterval, "")
	flag.IntVar(&interval, "i", DefaultInterval, "")

	flag.BoolVar(&listenAny, "listen-any", false, "")
	flag.BoolVar(&listenAny, "a", false, "")

	flag.StringVar(&listenPort, "listen-port", "", "")
	flag.StringVar(&listenPort, "p", "", "")

	flag.BoolVar(&showSender, "show-sender", false, "")
	flag.BoolVar(&showSender, "s", false, "")

	flag.BoolVar(&showTimestamp, "show-timestamp", false, "")
	flag.BoolVar(&showTimestamp, "t", false, "")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: udp-receiver [OPTION]...
Options:
  -a, -listen-any       : listen all available IP addresses (i.e. 0.0.0.0)
  -b, -buffer-size=NUM  : buffer size (default: %d)
  -p, -listen-port=NUM  : listen port number (required)
  -s, -show-sender      : show sender ([address:port])
  -t, -show-timestamp   : show timestamp

Experimental options:
  -i, -interval         : output interval (millisecond) (default %d)
`, DefaultBufferSize, DefaultInterval)
	}

	flag.Parse()

	if listenPort == "" {
		fmt.Fprintln(os.Stderr, "udp-viewer: ERROR: Please specify a port number (e.g. -p=4126)")
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

	for {
		n, remoteEP, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Fprintln(os.Stderr, "udp-viewer: ERROR:", err)
		}
		if showTimestamp {
			fmt.Printf(time.Now().Format("2006-01-02 15:04:05.00 "))
		}
		if showSender {
			fmt.Printf("[%s] ", remoteEP)
		}
		fmt.Printf("%s\n", string(buf[0:n]))

		time.Sleep(time.Duration(interval) * time.Millisecond) // FIXME
	}
}
