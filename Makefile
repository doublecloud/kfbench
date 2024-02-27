
compile:
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o target/kfbench-$(GOOS)-$(GOARCH) ./cmd/main.go
	cp target/kfbench-$(GOOS)-$(GOARCH) images/cache/kfbench-$(GOOS)-$(GOARCH)

linux-amd linux-amd64:
	GOOS=linux GOARCH=amd64 $(MAKE) compile

linux-arm linux-arm64:
	GOOS=linux GOARCH=arm64 $(MAKE) compile

build:
	$(MAKE) linux-amd64
	$(MAKE) linux-arm64

run: 
	./target/kfbench-linux-amd64 franz --brokers localhost:9092 --topic test --static-record --compression none --concurrency 1