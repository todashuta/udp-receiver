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
		interval      time.Duration
		listenAny     bool
		listenPort    string
		quote         bool
		showTimestamp bool
		showSender    bool
	)

	flag.IntVar(&bufSize, "buffer-size", DefaultBufferSize, "")
	flag.IntVar(&bufSize, "b", DefaultBufferSize, "")

	flag.DurationVar(&interval, "interval", DefaultInterval, "")
	flag.DurationVar(&interval, "i", DefaultInterval, "")

	flag.BoolVar(&listenAny, "listen-any", false, "")
	flag.BoolVar(&listenAny, "a", false, "")

	flag.StringVar(&listenPort, "listen-port", "", "")
	flag.StringVar(&listenPort, "p", "", "")

	flag.BoolVar(&quote, "q", false, "")

	flag.BoolVar(&showSender, "show-sender", false, "")
	flag.BoolVar(&showSender, "s", false, "")

	flag.BoolVar(&showTimestamp, "show-timestamp", false, "")
	flag.BoolVar(&showTimestamp, "t", false, "")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: udp-receiver [OPTION]...
Options:
  -a, -listen-any       : listen all available IP addresses (i.e. 0.0.0.0)
  -b, -buffer-size=NUM  : buffer size (default: %d)
  -i, -interval=STRING  : output interval (e.g. 2h45m, 1m, 2s, 300ms) (default %d)
  -p, -listen-port=NUM  : listen port number (required)
  -s, -show-sender      : show sender ([address:port])
  -t, -show-timestamp   : show timestamp
`, DefaultBufferSize, DefaultInterval)
	}

	flag.Parse()

	if listenPort == "" {
		fmt.Fprintln(os.Stderr, "udp-receiver: ERROR: Please specify a port number (e.g. -p=4126)")
		os.Exit(1)
	}

	var listenIP string
	if listenAny {
		listenIP = "0.0.0.0"
	} else {
		listenIP = "127.0.0.1"
	}

	var useInterval bool
	useInterval = interval > 0

	localEP, err := net.ResolveUDPAddr("udp", listenIP+":"+listenPort)
	if err != nil {
		fmt.Fprintln(os.Stderr, "udp-receiver: ERROR:", err)
		os.Exit(1)
	}

	conn, err := net.ListenUDP("udp", localEP)
	if err != nil {
		fmt.Fprintln(os.Stderr, "udp-receiver: ERROR:", err)
		os.Exit(1)
	}
	defer conn.Close()

	var tick <-chan time.Time
	if useInterval {
		tick = time.Tick(interval)
	}

	var format string
	if quote {
		format = "%q"
	} else {
		format = "%s"
	}

	ch := make(chan string, 2048)
	go func() {
		buf := make([]byte, bufSize)
		for {
			n, remoteEP, err := conn.ReadFromUDP(buf)
			if err != nil {
				fmt.Fprintln(os.Stderr, "udp-receiver: ERROR:", err)
				continue
			}

			var s string
			if showTimestamp {
				s += time.Now().Format("2006-01-02 15:04:05.00 ")
			}
			if showSender {
				s += fmt.Sprintf("[%s] ", remoteEP)
			}
			s += fmt.Sprintf(format, string(buf[0:n]))
			ch <- s
		}
	}()

	if useInterval {
		for {
			select {
			case <-tick:
				fmt.Println(<-ch)
			default:
				<-ch
			}
		}
	} else {
		for {
			fmt.Println(<-ch)
		}
	}
}
