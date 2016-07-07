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

## TODO

* harden the code, it's just a prototype
  * trickle read to protect from timeout
  * remove race condition
  * add test
* Make it production ready
  * add logging
  * add a pidfile
* Make it configurable
  * Limit the count of upstream connections
