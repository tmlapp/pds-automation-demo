FROM golang:1.19.1

LABEL maintainer="Eric Shanks <eshanks@purestorage.com>"

RUN mkdir /consumer

ADD . /consumer

WORKDIR /consumer

RUN go get . && go build -o main .

CMD ["./main"]