FROM ubuntu:bionic-20230530 as kfbench
LABEL maintainer="Misha Epikhin epikhinm@double.cloud"

ARG JDK=openjdk-11
# install openjdk from official repository
RUN mkdir -p /opt/kafka && \
    apt update && \
    apt install -y ${JDK}-jdk-headless
WORKDIR /opt/kafka

ARG KAFKA_VERSION=3.6.1
COPY cache/kafka_2.13-${KAFKA_VERSION}.tgz /opt/kafka

RUN bash -c "tar --strip-components=1 -zxvf kafka_2.13-$KAFKA_VERSION.tgz && rm -f kafka_2.13-$KAFKA_VERSION.tgz"
COPY server.properties /opt/kafka/config/kraft/server.properties
ENV PATH /opt/kafka/bin:$PATH


EXPOSE 9092/tcp
CMD bash -c 'export KAFKA_CLUSTER_ID=$(bin/kafka-storage.sh random-uuid) && \
    bin/kafka-storage.sh format -t $KAFKA_CLUSTER_ID -c config/kraft/server.properties && \
    bin/kafka-server-start.sh config/kraft/server.properties'
    # kafka-producer-perf-test.sh
