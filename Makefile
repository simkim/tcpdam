install: deps
	go install github.com/simkim/tcpdam
	go install github.com/simkim/tcpdam/cmd/tcpdam

deps:
	go get -d github.com/simkim/tcpdam
	go get -d github.com/simkim/tcpdam/cmd/tcpdam

build-bin:
	docker run -v `pwd`/build:/go/bin -v `pwd`:/go/src/github.com/simkim/tcpdam --rm golang go get github.com/simkim/tcpdam/...

build-image: build-bin
	docker build -t tcpdam .

dockerhub: build-image
	docker tag tcpdam simkim/tcpdam:latest
	docker push simkim/tcpdam:latest
