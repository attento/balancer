Balancer is a load-balancer for the Container age, easy and battery
included.

```
 HTTP[s] requests
   | |
   v v
[balancer]<--[ kvstore/file/swarm ]  // on-watch configuration
      |\___________________
      v                    v
web-serverA:portX  web-serverB:portY // upstreams and fallback
```

* [Home](/docs/index.md)
* source
    * [from config](/docs/config.md)
    * [from api](/docs/api.md)
* [keyworlds](/docs/keyworlds.md)
