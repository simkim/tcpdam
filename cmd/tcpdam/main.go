package main

import (
	"flag"
	"os"
	"strconv"

	"github.com/op/go-logging"
	"github.com/simkim/tcpdam"
)

var (
	localAddr        = flag.String("l", ":9999", "local address")
	remoteAddr       = flag.String("r", "127.0.0.1:80", "remote address")
	maxParkedProxies = flag.Int("max-parked", 0, "maximum parked connections")
	verbose          = flag.Bool("v", false, "show major events like open/close")
	debug            = flag.Bool("d", false, "show all debug events")
	pidFile          = flag.String("p", "", "pid file")
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

func teardownPidfile() error {
	err := os.Remove(*pidFile)
	if err != nil {
		log.Errorf("Can't remove pidfile : %s", err.Error())
	}
	return err
}

func setupPidfile() (error, bool) {
	if *pidFile == "" {
		return nil, false
	}
	file, err := os.OpenFile(*pidFile, os.O_CREATE|os.O_EXCL|os.O_WRONLY|os.O_TRUNC, 0777)
	if err != nil {
		return err, false
	}
	pid := os.Getpid()
	log.Debugf("Write pid %d to pidfile %s\n", pid, *pidFile)
	file.Write([]byte(strconv.Itoa(pid)))
	file.Close()
	return nil, true
}

var log = logging.MustGetLogger("tcpdam")

func main() {
	flag.Parse()
	setupLogging()
	err, hasPid := setupPidfile()
	if err != nil {
		log.Errorf("Can't create pid file : %s\n", err.Error())
		os.Exit(1)
	} else {
		if hasPid {
			defer teardownPidfile()
		}
	}
	log.Notice("tcpdam started")
	dam := tcpdam.NewDam(localAddr, remoteAddr)
	dam.Logger = log
	dam.MaxParkedProxies = *maxParkedProxies
	dam.Start()
	log.Notice("tcpdam stopped")
}
