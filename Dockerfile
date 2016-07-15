FROM golang:1.6.2
MAINTAINER Kevin Cantwell <kevin.cantwell@gmail.com>

# Add the project into the GOPATH
ADD . /go/src/github.com/kevin-cantwell/gotalks
WORKDIR /go/src/github.com/kevin-cantwell/gotalks

# Install gotalk
ENV GOBIN /go/bin
RUN go install ./...

# Install the present server
RUN go get -u golang.org/x/tools/cmd/present

# Usage: docker run -it -p 3999:3998 kevincantwell/gotalks:latest $(docker-machine ip) 3998
ENTRYPOINT ["./bin/start.sh"]
EXPOSE 3998