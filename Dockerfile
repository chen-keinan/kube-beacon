# Use an official golang runtime as a parent image
FROM golang:1.15-alpine as builder

ENV GO111MODULE=on

ADD . /src

WORKDIR /src/cmd/kube

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o kube-beacon .

FROM alpine:latest

RUN apk --no-cache add ca-certificates curl

WORKDIR /root/

COPY --from=builder /src/cmd/kube/kube-beacon .

RUN curl -L -o /usr/local/bin/kubectl "https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl"

RUN chmod +x /usr/local/bin/kubectl

ENV PATH "$PATH:/usr/local/bin/kubectl"

CMD ["./kube-beacon"]