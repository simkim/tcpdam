all: deps
	go install github.com/simkim/tcpdam
	go install github.com/simkim/tcpdam/cmd/tcpdam
deps:
	go get -d github.com/simkim/tcpdam
	go get -d github.com/simkim/tcpdam/cmd/tcpdam

docker:
	docker run -v `pwd`/build:/go/bin -v `pwd`:/go/src/github.com/simkim/tcpdam --rm golang go get github.com/simkim/tcpdam/...
	docker build -t tcpdam .

dockerhub: docker
	docker tag tcpdam simkim/tcpdam:latest
	docker push simkim/tcpdam:latest
