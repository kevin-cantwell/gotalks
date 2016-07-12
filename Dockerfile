FROM golang:1.6.2
MAINTAINER Kevin Cantwell <kevin.cantwell@gmail.com>

RUN go get -u golang.org/x/tools/cmd/present

ADD bin/start.sh /usr/bin/start
WORKDIR /root
EXPOSE 3999
ENTRYPOINT ["start"]