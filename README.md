balancer
========

Balancer is a load-balancer for the Container age.

## Features:

- Rest API to handle configuration.
- Is Crazy fast.
- Config watcher: push notification when the config file is changed, or the config on Consul/Etcd/Zookeeper is changed.
- Work with Docker and Docker-Swarm (*)
  automatically get nodes on services with Consul/Etcd/Zookeeper and docker-swarm
- If an Upstream is down or give bad response, it use the fallback. (*)

* to be completed.

```
 HTTP[s] requests 
   | |
   v v
[balancer]<--[ kvstore/file/swarm ]  // on-watch configuration
      |\___________________
      v                    v
web-serverA:portX  web-serverB:portY // upstreams and fallback
```

## Configuration is a list of Servers

Servers are available on port and filters

```
{  
   "address":":80",
   "filter":{                      // optional used to filter requests
      "Hosts":["www.website.com"], // optional used to filter on hosts
      "Schemes":["http"],          
      "PathPrefix":""              // optional used to filter some routes
   },
   "upstreams":{  
      [{ "Target":"127.0.0.1",
         "Port":80,
         "Priority":1,             // optional
         "Weight":2                // optional
      }]
   }
}
```

## Use with the Dockerfile

Build

`docker build -t balancer .`

Run

`docker run -p 8000:80 -p 9123:9123  balancer`


## API

### Config `GET /`

Return the whole configuration.

### Servers 

#### `GET /server/{address}` 

Return the Server Object

#### Upstream `GET /server/:address/upstream/{target}-{port}` 

Return the Upstream object

#### Upstream `POST|PUT /server/{address}/upstream/{target}-{port}` 

Create a new Object Stream or Updated it

Request body: `{"Target":"127.0.0.1","Port":8080,"Priority":1,"Weight":2}`

#### Filter `GET /server/{address}/filter`

Get the Server Filter

#### Filter `POST|PUT /server/{address}/filter`

Create or Update the filter

Request body `{"Hosts":null,"Schemes":["",""],"PathPrefix":""}`

## Test it

```
go get ./...
go test ./...
```
