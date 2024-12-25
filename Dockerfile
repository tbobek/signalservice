FROM golang:1.20 AS builder
WORKDIR /go/src/github.com/tbobek/signalserver
COPY *.go ./
COPY go.mod ./ 
RUN go get ./... 
# RUN GO111MODULE="on" CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app
RUN go build -a -ldflags "-linkmode external -extldflags '-static' -s -w" -o app .
FROM scratch
ENV PORT ${PORT_SIGNALSERVICE}
WORKDIR /root/
COPY --from=builder /go/src/github.com/tbobek/signalserver/app ./
CMD ["./app"]
EXPOSE ${PORT}
