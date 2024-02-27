
VERSION=3.6.1

.PHONY=clean
clean:
	rm -rf cache

cache:
	mkdir -p cache
	curl -L -o cache/kafka_2.13-3.6.1.tgz "https://dlcdn.apache.org/kafka/3.6.1/kafka_2.13-3.6.1.tgz"

# 3.6.1-11.0.22-amzn:
# 	docker build --build-arg JDK=11.0.22-amzn --build-arg KAFKA_VERSION=3.6.1 -t doublecloud/kfbench:3.6.1-11.0.22-amzn .

# 3.6.1-17.0.10-amzn:
# 	docker build --build-arg JDK=17.0.10-amzn --build-arg KAFKA_VERSION=3.6.1 -t doublecloud/kfbench:3.6.1-17.0.10-amzn .

# 3.6.1-21.0.2-amzn:
# 	docker build --build-arg JDK=21.0.2-amzn --build-arg KAFKA_VERSION=3.6.1 -t doublecloud/kfbench:3.6.1-21.0.2-amzn .

# 3.6.1-21.0.2


amzn:
	docker build --build-arg JDK=21.0.2-amzn  --build-arg KAFKA_VERSION=3.6.1 -t doublecloud/kfbench:3.6.1-21.0.2-amzn -f jammy.Dockerfile
	docker build --build-arg JDK=17.0.10-amzn --build-arg KAFKA_VERSION=3.6.1 -t doublecloud/kfbench:3.6.1-17.0.10-amzn -f jammy.Dockerfile
	docker build --build-arg JDK=11.0.22-amzn --build-arg KAFKA_VERSION=3.6.1 -t doublecloud/kfbench:3.6.1-11.0.22-amzn -f jammy.Dockerfile
	docker build --build-arg JDK=8.0.402-amzn --build-arg KAFKA_VERSION=3.6.1 -t doublecloud/kfbench:3.6.1-8.0.402-amzn -f jammy.Dockerfile

bionic:
	docker build --platform linux-arm64 --build-arg JDK=openjdk-11  --build-arg KAFKA_VERSION=3.6.1 -t doublecloud/kfbench:3.6.1-openjdk-11 -f bionic.Dockerfile
	docker build --platform linux-arm64 --build-arg JDK=openjdk-17  --build-arg KAFKA_VERSION=3.6.1 -t doublecloud/kfbench:3.6.1-openjdk-17 -f bionic.Dockerfile

# albba:
# 	echo $(MAKE) $(VERSION)-21.0.2-amzn
# 	echo $(MAKE) $(VERSION)-17.0.10-amzn
# 	echo $(MAKE) $(VERSION)-11.0.22-amzn
# 	echo $(MAKE) $(VERSION)-8.0.402-amzn
# build:
	