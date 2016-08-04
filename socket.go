package tcpdam

import (
	"net"
	"os"
	"strings"
)

func (dam *Dam) executeCommand(args []string) {
	switch args[0] {
	case "open":
		dam.Open()
	case "close":
		dam.Close()
	case "set-remote":
		if len(args) > 1 {
			dam.SetRemoteAddr(args[1])
		} else {
			dam.Logger.Errorf("Missing argument for %s", args[0])
		}
	default:
		dam.Logger.Errorf("Command not found: %s", args[0])
	}
}

func (dam *Dam) StartControlSocket(path string) {
	dam.Logger.Debugf("Starting control socket at %s", path)
	l, err := net.ListenUnix("unix", &net.UnixAddr{Name: path, Net: "unix"})
	if err != nil {
		panic(err)
	}
	defer os.Remove(path)

	for {
		conn, err := l.AcceptUnix()
		if err != nil {
			panic(err)
		}
		var buf [1024]byte
		n, err := conn.Read(buf[:])
		if err != nil {
			panic(err)
		}
		args := strings.Split(strings.TrimSpace(string(buf[:n])), " ")
		dam.executeCommand(args)
		conn.Close()
	}
}

func SendControlCommand(path string, cmd, arg []string) {
	conn, err := net.DialUnix("unix", nil,
		&net.UnixAddr{Name: path, Net: "unix"})
	if err != nil {
		panic(err)
	}

	_, err = conn.Write([]byte("hello"))
	if err != nil {
		panic(err)
	}
	conn.Close()
}
