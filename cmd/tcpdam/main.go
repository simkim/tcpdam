package main

import (
	"flag"
	"os"
	"strconv"

	"github.com/op/go-logging"
	"github.com/simkim/tcpdam"
)

func configFromEnv(key string, _default string) string {
	if os.Getenv(key) != "" {
		return os.Getenv(key)
	} else {
		return _default
	}
}

func configFromEnvBool(key string, _default bool) bool {
	if os.Getenv(key) != "" {
		rc, err := strconv.ParseBool(os.Getenv(key))
		if err != nil {
			log.Warning("Invalid value %s in env var %s", os.Getenv(key), key)
			return _default
		}
		return rc
	} else {
		return _default
	}
}

var (
	listenAddr       = flag.String("l", configFromEnv("TCPDAM_LISTEN_ADDRESS", ":9999"), "listen address (TCPDAM_LISTEN_ADDRESS)")
	remoteAddr       = flag.String("r", configFromEnv("TCPDAM_REMOTE_ADDRESS", "127.0.0.1:80"), "remote address (TCPDAM_REMOTE_ADDRESS)")
	maxParkedProxies = flag.Int("max-parked", 0, "maximum parked connections")
	verbose          = flag.Bool("v", configFromEnvBool("TCPDAM_VERBOSE", false), "show major events like open/close (TCPDAM_VERBOSE)")
	debug            = flag.Bool("d", configFromEnvBool("TCPDAM_DEBUG", false), "show all debug events (TCPDAM_DEBUG)")
	pidFile          = flag.String("p", configFromEnv("TCPDAM_PIDFILE", ""), "pid file (TCPDAM_PIDFILE)")
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
	log.Noticef("tcpdam started (%s -> %s)", *listenAddr, *remoteAddr)
	dam := tcpdam.NewDam(listenAddr, remoteAddr)
	dam.Logger = log
	dam.MaxParkedProxies = *maxParkedProxies
	dam.Start()
	log.Notice("tcpdam stopped")
}
