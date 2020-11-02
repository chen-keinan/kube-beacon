# Use an official golang runtime as a parent image
FROM golang:1.14

# Set the working directory to /app
RUN  apt-get update -y

RUN  apt-get install -y gawk

RUN mkdir -p /app

ADD beacon /app

CMD ["/app/beacon"]