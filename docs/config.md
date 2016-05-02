# Configuration

It supports configuration from file located in
`/etc/balancer.json`.  Create this file and
start balancer with the `file` exporter.

```
balancer run -e file
```

## Configuration is a list of Servers
Servers are available on port and filters
```
[
    {
       "address":":80",
       "filter":{                      // optional used to filter requests
          "Hosts":["www.website.com"], // optional used to filter on hosts
          "Schemes":["http"],
          "PathPrefix":""              // optional used to filter some routes
       },
       "upstreams":{
         [
           {
             "Target":"127.0.0.1",
             "Port":80,
             "Priority":1,             // optional
             "Weight":2                // optional
           },
           {
             "Target":"127.0.0.1",
             "Port":8080
           }
         ]
       }
    }
]
```

* [Home](/docs/index.md)
* source
    * [from config](/docs/config.md)
    * [from api](/docs/api.md)
* [keyworlds](/docs/keyworlds.md)
