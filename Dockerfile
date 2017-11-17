FROM golang:1.9.2-alpine3.6

LABEL maintainer="support@inwecrypto.com"

COPY . /go/src/github.com/inwecrypto/neo-insight

RUN go install github.com/inwecrypto/neo-insight && rm -rf /go/src

VOLUME ["/etc/inwecrypto/insight/neo"]

WORKDIR /etc/inwecrypto/insight/neo

CMD ["/go/bin/neo","--conf","/etc/inwecrypto/insight/neo/insight.json"]