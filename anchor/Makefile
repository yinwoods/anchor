BINARY_NAME = anchor
SCHEDULER_NAME = scheduler

build:
	go build -o $(BINARY_NAME) -v

run: 
	go build -o $(BINARY_NAME) -v 
	./$(BINARY_NAME)

scheduler:
	go build -o scheduler/$(SCHEDULER_NAME) -v

test:
	go test -v .

linux-build:
	GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME) -v

darwin-build:
	GOOS=darwin GOARCH=amd64 go build -o $(BINARY_NAME) -v

windows-build:
	GOOS=windows GOARCH=amd64 go build -o $(BINARY_NAME).exe -v

docker:
	docker build -t $(BINARY_NAME) .

clean:
	go clean
	rm -f $(BINARY_NAME).exe

.PHONY: build run test linux-build darwin-build windows-build docker clean
