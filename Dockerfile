FROM golang:1.8

RUN mkdir -p /go/src/github.com/wolakec/makhzen
WORKDIR /go/src/github.com/wolakec/makhzen

COPY . .

RUN go get -v -v ./...
RUN go install -v ./...

CMD ["makhzen"]

EXPOSE 5000