FROM golang:1.4.1-onbuild


WORKDIR /go/src/github.com/tirami/tfw-application-server
ADD . /go/src/github.com/tirami/tfw-application-server/

# go get all of the dependencies

#RUN go get github.com/tirami/tfw-application-server

# install the app
RUN go install github.com/tirami/tfw-application-server

EXPOSE 8080

CMD []
ENTRYPOINT ["/go/bin/tfw-application-server"]