FROM golang:1.3-cross
ADD . /go/src/github.com/t0mk/bitbucket_jenkins_proxy
WORKDIR /go/src/github.com/t0mk/bitbucket_jenkins_proxy
ENV GOOS linux
ENV GOARCH amd64
RUN go get
ENTRYPOINT ["/go/src/github.com/t0mk/bitbucket_jenkins_proxy/make.sh"]
