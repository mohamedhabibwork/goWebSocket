FROM golang:1.16-alpine as builder

RUN mkdir /socket
ADD . /socket
WORKDIR /socket
RUN go clean --modcache
COPY src/go.mod ./
COPY src/go.sum ./
RUN go mod download
#RUN CGO_ENABLE=0 GOOS=linux go build -a -installsuffix cgo -o main .

#FORM alpine:latest
#RUN apk --no-cache add ca-certificates
#RUN apk add --no-cache git make musl-dev go

COPY ./*.go ./src

# Build
RUN go build -o /socket

# This is for documentation purposes only.
# To actually open the port, runtime parameters
# must be supplied to the docker command.
EXPOSE 8080

# (Optional) environment variable that our dockerised
# application can make use of. The value of environment
# variables can also be set via parameters supplied
# to the docker command on the command line.
ENV PORT=8080
# Run
#USER root
#RUN chmod +x *
CMD [ "/socket" ]
