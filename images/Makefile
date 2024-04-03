VERSION=3.6.1
IMAGE=kfbench
REGISTRY=?

.PHONY=clean
clean:
	rm -rf cache

cache:
	mkdir -p cache
	curl -L -o cache/kafka_2.13-$(VERSION).tgz "https://dlcdn.apache.org/kafka/$(VERSION)/kafka_2.13-$(VERSION).tgz"

amzn: 21.0.2-amzn 17.0.10-amzn 11.0.22-amzn 8.0.402-amzn
%-amzn:
	podman build \
		--platform linux/amd64 \
		--build-arg JDK=$@ \
		--build-arg KAFKA_VERSION=$(VERSION) \
		-t $(IMAGE):$(VERSION)-$@ \
		-f jammy.Dockerfile
	podman tag $(IMAGE):$(VERSION)-$@ $(REGISTRY)/$(IMAGE):$(VERSION)-$@-amd64
	podman push $(REGISTRY)/$(IMAGE):$(VERSION)-$@
	podman build \
		--platform linux/arm64 \
		--build-arg JDK=$@ \
		--build-arg KAFKA_VERSION=$(VERSION) \
		-t $(IMAGE):$(VERSION)-$@ \
		-f jammy.Dockerfile
	podman tag $(IMAGE):$(VERSION)-$@ $(REGISTRY)/$(IMAGE):$(VERSION)-$@-arm64
	podman push $(REGISTRY)/$(IMAGE):$(VERSION)-$@

# Test outdated environments for historical comparisions
bionic:
	docker build --build-arg JDK=openjdk-11  --build-arg KAFKA_VERSION=3.6.1 -t doublecloud/kfbench:3.6.1-openjdk-11 -f bionic.Dockerfile
	docker build --build-arg JDK=openjdk-17  --build-arg KAFKA_VERSION=3.6.1 -t doublecloud/kfbench:3.6.1-openjdk-17 -f bionic.Dockerfile
