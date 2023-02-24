FROM ubuntu

LABEL maintainer="Eric Shanks <eshanks@purestorage.com>"

RUN apt update && apt install vim default-jre wget curl -y

RUN wget https://archive.apache.org/dist/kafka/3.0.0/kafka_2.13-3.0.0.tgz && tar xzf kafka_2.13-3.0.0.tgz

