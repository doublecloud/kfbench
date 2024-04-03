FROM ubuntu:jammy-20240227 as kfbench
LABEL maintainer="Misha Epikhin epikhinm@double.cloud"

# install sdkman for different jvms
ENV SDKMAN_INIT "/usr/local/sdkman/bin/sdkman-init.sh"
ENV SDKMAN_DIR="/usr/local/sdkman"
RUN mkdir -p /opt/kafka && \
    apt update && \
    apt install -y curl unzip zip && \
    curl -s "https://get.sdkman.io" | bash && \
    /bin/bash -c "source ${SDKMAN_INIT} && sdk version && sdk update && sdk upgrade"
COPY 50-sdkman.sh /etc/profile.d/50-sdkman.sh

ARG JDK=11.0.22-amzn
ARG KAFKA_VERSION=3.6.1
RUN bash -c "source ${SDKMAN_INIT} && \
    sdk install java ${JDK} && \
    sdk default java ${JDK} && \
    sdk flush archives && \
    sdk flush temp \
    chmod -R 755 ${SDKMAN_DIR}"
ENV PATH ${SDKMAN_DIR}/candidates/java/current/bin:$PATH
ENV JAVA_HOME=${SDKMAN_DIR}/candidates/java/current/
WORKDIR /opt/kafka

ARG KAFKA_VERSION
COPY cache/kafka_2.13-${KAFKA_VERSION}.tgz /opt/kafka

RUN bash -c "tar --strip-components=1 -zxvf kafka_2.13-$KAFKA_VERSION.tgz && rm -f kafka_2.13-$KAFKA_VERSION.tgz"
COPY server.properties /opt/kafka/config/kraft/server.properties
ENV PATH /opt/kafka/bin:$PATH

EXPOSE 9092/tcp
CMD bash -c 'export KAFKA_CLUSTER_ID=$(bin/kafka-storage.sh random-uuid) && \
    bin/kafka-storage.sh format -t $KAFKA_CLUSTER_ID -c config/kraft/server.properties && \
    bin/kafka-server-start.sh config/kraft/server.properties'
