FROM golang:1-alpine

COPY . $GOPATH/src/github.com/attento/balancer
WORKDIR $GOPATH/src/github.com/attento/balancer

RUN apk update && \
   apk add git

RUN  go get ./... && \
     go install -v

CMD ["balancer", "r"]