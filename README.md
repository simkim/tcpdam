# tcpdam

## Overview

**tcpdam** is a parking proxy for your tcp connection. When your upstream server is ready, send SIGUSR1 to flush connection to the server.

## Compilation

    make

## Usage

    tcpdam -l LISTEN-HOST:LISTEN-PORT -r REMOTE-HOST:REMOTE-PORT

will setup a listening tcp server, when something connect it is parked in a waiting list

    killall -USR1 tcpdam

will lookup REMOTE-HOST once, and open the dam : the parked and new connections are proxified to REMOTE-HOST.

    killall -USR2 tcpdam

will re-close the dam and start to park new connections,

## Configuration
| Env | Command line | Description |
| --- | ------------ | ----------- |
| TCPDAM_LISTEN_ADDRESS | -l | Listen address for incoming connections |
| TCPDAM_REMOTE_ADDRESS | -h | Remote address where the connections will be flushed |
| TCPDAM_DEBUG | -d | Show all information |
| TCPDAM_VERBOSE | -v | Show some information
| TCPDAM_PIDFILE | -p | File which will contain the pid of the dam |
| TCPDAM_CTRLSOCKET | -ctrl-socket | Unix socket to control the dam |
| TCPDAM_MAX_FLUSHING | -max-flushing | Max number of open remote connections |
| TCPDAM_MAX_PARKED | -max-parked | Max number of connections in the queue |
| TCPDAM_OPEN | -open | Start the dam open |
|   | -c | command to send to a running dam |

### Remote Commands

| Command | Description |
| ------- | ----------- |
| open    | open the dam unless already open |
| close   | close the dam unless already closed |
| set-remote HOST:PORT | switch the remote address, will be used at the next open |

## Limits

### Open file descriptors

You need to configure your maximum number of file descriptor to allow a high number of parked connections

 * debian  : configure /etc/security/limits.conf
 * docker  : docker run --ulimit nofile=100000:100000
 * compose : see docker-compose.yml

### netfilter conntrack size

TODO

### tcp TIME_WAIT

TODO

## Docker

### From the hub

To run the dam

    docker run -p 9999:9999 --rm -ti --name tcpdam_test  simkim/tcpdam tcpdam -r google.com:80

To open the dam

    docker exec tcpdam_test killall -USR1 tcpdam

To open the dam, wait for connections to terminate and quit

    docker stop tcpdam_test

### From local build

    Edit docker-compose.yml

    docker run -v `pwd`/build:/go/bin -v `pwd`:/go/src/github.com/simkim/tcpdam --rm golang go get github.com/simkim/tcpdam/...
    docker-compose build
    docker-compose up
