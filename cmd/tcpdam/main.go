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
	}
	return _default
}

func configFromEnvInt(key string, _default int) int {
	if os.Getenv(key) != "" {
		rc, err := strconv.Atoi(os.Getenv(key))
		if err != nil {
			log.Warning("Invalid value %s in env var %s", os.Getenv(key), key)
			return _default
		}
		return rc
	}
	return _default
}

func configFromEnvBool(key string, _default bool) bool {
	if os.Getenv(key) != "" {
		rc, err := strconv.ParseBool(os.Getenv(key))
		if err != nil {
			log.Warning("Invalid value %s in env var %s", os.Getenv(key), key)
			return _default
		}
		return rc
	}
	return _default
}

var (
	listenAddr         = flag.String("l", configFromEnv("TCPDAM_LISTEN_ADDRESS", ":9999"), "listen address (TCPDAM_LISTEN_ADDRESS)")
	remoteAddr         = flag.String("r", configFromEnv("TCPDAM_REMOTE_ADDRESS", "127.0.0.1:80"), "remote address (TCPDAM_REMOTE_ADDRESS)")
	maxParkedProxies   = flag.Int("max-parked", configFromEnvInt("TCPDAM_MAX_PARKED", 100000), "maximum parked connections")
	maxFlushingProxies = flag.Int("max-flushing", configFromEnvInt("TCPDAM_MAX_FLUSHING", 10), "maximum flushing connections")
	verbose            = flag.Bool("v", configFromEnvBool("TCPDAM_VERBOSE", false), "show major events like open/close (TCPDAM_VERBOSE)")
	debug              = flag.Bool("d", configFromEnvBool("TCPDAM_DEBUG", false), "show all debug events (TCPDAM_DEBUG)")
	pidFile            = flag.String("p", configFromEnv("TCPDAM_PIDFILE", ""), "pid file (TCPDAM_PIDFILE)")
	open               = flag.Bool("open", configFromEnvBool("TCPDAM_OPEN", false), "start already open (TCPDAM_OPEN)")
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

func setupPidfile() (bool, error) {
	if *pidFile == "" {
		return false, nil
	}
	file, err := os.OpenFile(*pidFile, os.O_CREATE|os.O_EXCL|os.O_WRONLY|os.O_TRUNC, 0777)
	if err != nil {
		return false, err
	}
	pid := os.Getpid()
	log.Debugf("Write pid %d to pidfile %s\n", pid, *pidFile)
	file.Write([]byte(strconv.Itoa(pid)))
	file.Close()
	return true, nil
}

var log = logging.MustGetLogger("tcpdam")

func main() {
	flag.Parse()
	setupLogging()
	hasPid, err := setupPidfile()
	if err != nil {
		log.Errorf("Can't create pid file : %s\n", err.Error())
		os.Exit(1)
	} else {
		if hasPid {
			defer teardownPidfile()
		}
	}
	log.Noticef("tcpdam started (%s -> %s)", *listenAddr, *remoteAddr)
	dam := tcpdam.NewDam(*listenAddr, *remoteAddr, *maxParkedProxies, *maxFlushingProxies)
	dam.Logger = log
	err = dam.Start(*open)
	if err != nil {
		log.Errorf("An error occured: %s", err.Error())
	}
	log.Notice("tcpdam stopped")
}
