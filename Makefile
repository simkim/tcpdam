all: deps
	go install github.com/simkim/tcpdam
	go install github.com/simkim/tcpdam/cmd/tcpdam
deps:
	go get -d github.com/simkim/tcpdam
	go get -d github.com/simkim/tcpdam/cmd/tcpdam
docker:
	go build github.com/simkim/tcpdam
	go build -o build/tcpdam github.com/simkim/tcpdam/cmd/tcpdam
	docker build -t tcpdam .

dockerhub: docker
	docker tag tcpdam simkim/tcpdam:latest
	docker push simkim/tcpdam:latest
