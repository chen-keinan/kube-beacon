# Use an official golang runtime as a parent image
FROM golang:1.17-alpine as builder

ENV GO111MODULE=on

ADD . /src

WORKDIR /src/cmd/kube

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o kube-beacon .

FROM golang:1.15-alpine

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /src/cmd/kube/kube-beacon .

CMD ["./kube-beacon"]