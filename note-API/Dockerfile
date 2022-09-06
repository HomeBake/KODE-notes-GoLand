FROM golang:1.19

ADD . /go/src/note-API
WORKDIR /go/src/note-API
ENTRYPOINT ["/go/bin/note-API"]

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go get -d -v ./...
RUN go install -v ./...

RUN go build -o /docker-gs-ping

EXPOSE 8080

CMD [ "/docker-gs-ping" ]