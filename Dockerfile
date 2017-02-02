FROM golang
EXPOSE 9999

ADD . /go/src/github.com/simkim/tcpdam
RUN go get -v github.com/simkim/tcpdam/...

USER nobody
ENTRYPOINT [ "tcpdam" ]
