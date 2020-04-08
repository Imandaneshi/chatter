FROM golang:alpine

WORKDIR $GOPATH/src/github.com/imandaneshi/chatter

COPY . .

RUN go mod download

RUN go install

CMD ["chatter", "server"]