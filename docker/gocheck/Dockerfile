FROM ubuntu:18.04

ENV GO_DL_VERSION 1.15.5

RUN apt-get update && apt-get install -y vim wget monitoring-plugins

RUN wget "https://golang.org/dl/go${GO_DL_VERSION}.linux-amd64.tar.gz" \
    && tar -C /usr/local/ -xzf "go${GO_DL_VERSION}.linux-amd64.tar.gz"

ENV PATH $PATH:/usr/local/go/bin

WORKDIR /app

COPY src/ /app/

RUN go build

COPY testFile.hcl /app/etc/
