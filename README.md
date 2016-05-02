balancer
========

[![Build Status](https://travis-ci.org/attento/balancer.svg?branch=master)](https://travis-ci.org/attento/balancer) 

Balancer is a load-balancer for the Container age, easy and battery included.

## Features:

- Rest API to handle configuration.
- Is Crazy fast.
- Hot-reloading of configuration: filters and Upstreams. 
- Graceful shutdown http connections and drain it before close connections.
- Event system.

Coming soon ...

- HTTPS
- Config watcher: push notification when the config file is changed, or the config on Consul/Etcd/Zookeeper is changed.
- Work with Docker and Docker-Swarm 
  automatically get nodes from services like Consul/Etcd/Zookeeper and docker-swarm
- Fallback and Circuit breaker.
- Let's encrypt integration to privide automatic SSL services.
- WebSocket.


```
 HTTP[s] requests
   | |
   v v
[balancer]<--[ kvstore/file/swarm ]  // on-watch configuration
      |\___________________
      v                    v
web-serverA:portX  web-serverB:portY // upstreams and fallback
```

## Use with the Dockerfile
`docker build -t balancer . && docker run -p 8000:80 -p 9123:9123  balancer`

## Simple usage:

Then just use the API like `curl localhost:9123`

1. Store a new Upstream binding :8080 `curl -X
   POST  localhost:9123/server --data
   '{"address":":8080","upstreams":[{"Target":"www.google.com","Port":80}]}'`
2. View the server: `curl
   localhost:9123/server/:8080`
3. Use it! `curl  localhost:8080` The proxy
   should response with the google homepage...

## Test it

```
go get ./...
go test ./...
```
