FROM golang:1-alpine

COPY . $GOPATH/src/github.com/liuggio/balancer
WORKDIR $GOPATH/src/github.com/liuggio/balancer

RUN apk update && \
   apk add git

RUN  go get ./... && \
     go install -v

CMD ["balancer", "r"]