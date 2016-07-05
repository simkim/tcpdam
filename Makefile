all: deps
	go install github.com/simkim/tcpdam
	go install github.com/simkim/tcpdam/cmd/tcpdam
deps:
	go get -d github.com/simkim/tcpdam
	go get -d github.com/simkim/tcpdam/cmd/tcpdam
