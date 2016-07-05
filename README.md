# tcpdam

## Overview

**tcpdam** is a parking proxy for your tcp connection. When your upstream server is ready, send SIGUP to flush connection to the server.

## Compilation

    make

## Usage

    tcpdam -l LISTEN-HOST:LISTEN-PORT -r REMOTE-HOST:REMOTE-PORT

will setup a listening tcp server, when something connect it is parked in a waiting list

    killall -HUP tcpdam

will lookup REMOTE-HOST then establish and proxy connection to it.

## TODO

* harden the code, it's just a prototype
* add a pidfile
* trickle read to protect from timeout
* remove flushed proxy from waiting proxy
* add logging
* add test
