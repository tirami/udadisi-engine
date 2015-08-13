FROM golang:1.4.1-onbuild

RUN apt-get update && apt-get install -y ca-certificates git-core ssh
ADD keys/my_key_rsa /root/.ssh/id_rsa
RUN chmod 700 /root/.ssh/id_rsa
RUN echo "Host github.com\n\tStrictHostKeyChecking no\n" >> /root/.ssh/config
RUN git config --global url.ssh://git@github.com

WORKDIR /gopath/src/github.com/tirami/tfw-application-server
ADD . /gopath/src/github.com/tirami/tfw-application-server/

# go get all of the dependencies

RUN go get github.com/tirami/tfw-application-server

# set env variables
ENV POSTGRES_DB "localhost"

EXPOSE 8080

CMD []
ENTRYPOINT ["/gopath/bin/tfw-application-server"]