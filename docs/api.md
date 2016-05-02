## API

You can decide to start this service without any
configuration file and use the API to start your
servers and add your upstream.

Then just use the API like `curl localhost:9123`

1. Store a new Upstream binding :8080 `curl -X
   POST  localhost:9123/server --data
   '{"address":":8080","upstreams":[{"Target":"www.google.com","Port":80}]}'`
2. View the server: `curl
   localhost:9123/server/:8080`
3. Use it! `curl  localhost:8080` The proxy
   should response with the google site...

## Config

`GET /` Return the whole configuration.

## Servers

`GET /server/{address}` Return the Server Object
`DELETE /server/{address}`

## Upstream

- `GET /server/:address/upstream/{target}-{port}`  Return the Upstream object

- `POST|PUT /server/{address}/upstream/{target}-{port}` Create a new Object Stream or Updated it

Request body: `{"Target":"127.0.0.1","Port":8080,"Priority":1,"Weight":2}`

- `DELETE /server/{address}/upstream/{target}-{port}` Delete

## Filter 

- `GET /server/{address}/filter` Get the Server Filter

- `POST|PUT /server/{address}/filter` Create or Update the filter
Request body `{"Hosts":null,"Schemes":["",""],"PathPrefix":""}`

* [Home](/docs/index.md)
* source
    * [from config](/docs/config.md)
    * [from api](/docs/api.md)
* [keyworlds](/docs/keyworlds.md)
