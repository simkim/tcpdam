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

## Docker

### Running

To run the dam

    docker run -p 9999:9999 --r-ti --name hooverdam  simkim/tcpdam tcpdam -r hoover.com:80

To open the dam

    docker exec hooverdam killall -USR1 tcpdam

### Build your tcpdam image

    docker run -v `pwd`/build:/go/bin -v `pwd`:/go/src/github.com/simkim/tcpdam --rm golang go get github.com/simkim/tcpdam/...
    docker build -t tcpdam .

