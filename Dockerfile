FROM golang
EXPOSE 9999

RUN go get -v github.com/simkim/tcpdam/...

USER nobody
ENTRYPOINT [ "tcpdam" ]
