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

### Environment variables
* TCPDAM_LISTEN_ADDRESS
* TCPDAM_REMOTE_ADDRESS
* TCPDAM_DEBUG
* TCPDAM_VERBOSE
* TCPDAM_PIDFILE

### Command line switch
* -l
* -h
* -d
* -v
* -p
* -max-parked

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
