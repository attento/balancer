## Api
One of the balancer's goal is to be extensible, it is api first and you
can to all with that, use our core to do it awesome!

## Filter

# Keyworlds
A lit of important worlds to work with balancer.

## Server(s)
`:80`, `api.youtsite.int:8081`.

When you start a new server you are binding a new port to receive and
dispatch requests. It is defined from an address a port and it could have
a filter and much upstream.
balancer supports one or more servers in the meantime.

## Source(s)
You can decide to use any source and use only the api but balancer
support different dynamic way to auto-configure itself like file, consul,
swarm..

Ok I am a lier, at the moment we support only the one source, but we have
big plan!

## Upstream(s)
An upstream is the representation of a single node behind your server.
It is described from a Target (es. 127.0.0.10) a Port a Priority and a
Weight.
Each server has one or more upstreams.
