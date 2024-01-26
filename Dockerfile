FROM golang:1.21

RUN go install github.com/cespare/reflex@latest

WORKDIR $GOPATH/src/github.com/laouji/fizz

EXPOSE 5000

CMD $GOPATH/bin/fizz
