package main

import (
	"flag"
	"os"

	"github.com/op/go-logging"
	"github.com/simkim/tcpdam"
)

var (
	localAddr        = flag.String("l", ":9999", "local address")
	remoteAddr       = flag.String("r", "127.0.0.1:80", "remote address")
	maxParkedProxies = flag.Int("max-parked", 0, "maximum parked connections")
	verbose          = flag.Bool("v", false, "show major events like open/close")
	debug            = flag.Bool("d", false, "show all debug events")
)

func setupLogging() {
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	var format = logging.MustStringFormatter(
		`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
	)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	if !*debug {
		backendLeveled := logging.AddModuleLevel(backendFormatter)
		if *verbose {
			backendLeveled.SetLevel(logging.NOTICE, "")
		} else {
			backendLeveled.SetLevel(logging.WARNING, "")
		}
		logging.SetBackend(backendLeveled)
	} else {
		logging.SetBackend(backendFormatter)
	}

}

var log = logging.MustGetLogger("tcpdam")

func main() {
	flag.Parse()
	setupLogging()
	log.Notice("tcpdam started")
	dam := tcpdam.NewDam(localAddr, remoteAddr)
	dam.Logger = log
	dam.MaxParkedProxies = *maxParkedProxies
	dam.Start()
}
