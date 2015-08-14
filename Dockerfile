FROM golang:1.4.1-onbuild


WORKDIR /go/src/github.com/tirami/udadisi-engine
ADD . /go/src/github.com/tirami/udadisi-engine/

# go get all of the dependencies

# install the app
RUN go install github.com/tirami/udadisi-engine

EXPOSE 8080

CMD []
ENTRYPOINT ["/go/bin/udadisi-engine"]