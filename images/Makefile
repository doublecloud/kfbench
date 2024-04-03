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
        -t $(IMAGE):$(VERSION)-$@-amd64 \
        -f jammy.Dockerfile
    podman tag $(IMAGE):$(VERSION)-$@-amd64 $(REGISTRY)/$(IMAGE):$(VERSION)-$@-amd64
    podman push $(REGISTRY)/$(IMAGE):$(VERSION)-$@-amd64
    podman build \
        --platform linux/arm64 \
        --build-arg JDK=$@ \
        --build-arg KAFKA_VERSION=$(VERSION) \
        -t $(IMAGE):$(VERSION)-$@-arm64 \
        -f jammy.Dockerfile
    podman tag $(IMAGE):$(VERSION)-$@-arm64 $(REGISTRY)/$(IMAGE):$(VERSION)-$@-arm64
    podman push $(REGISTRY)/$(IMAGE):$(VERSION)-$@-arm64

# Test outdated environments for historical comparisions
openjdk: openjdk-17 openjdk-11
openjdk-%:
    podman build \
        --platform linux/amd64 \
        --build-arg JDK=$@ \
        --build-arg KAFKA_VERSION=$(VERSION) \
        -t $(IMAGE):$(VERSION)-$@-amd64 \
        -f bionic.Dockerfile
    podman tag $(IMAGE):$(VERSION)-$@-amd64 $(REGISTRY)/$(IMAGE):$(VERSION)-$@-amd64
    podman push $(REGISTRY)/$(IMAGE):$(VERSION)-$@-amd64
    podman build \
        --platform linux/arm64 \
        --build-arg JDK=$@ \
        --build-arg KAFKA_VERSION=$(VERSION) \
        -t $(IMAGE):$(VERSION)-$@-arm64 \
        -f bionic.Dockerfile
    podman tag $(IMAGE):$(VERSION)-$@-arm64 $(REGISTRY)/$(IMAGE):$(VERSION)-$@-arm64
    podman push $(REGISTRY)/$(IMAGE):$(VERSION)-$@-arm64
