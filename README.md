# Juggler

Reverse proxy that controls the traffic and help replace the old HTTP service with the new one.

Juggler forwards the request against two backends, returns the response to the client and emits the diff between the two responses.

## Features
- Transparent to the client
- Driven by configuration
- Small footprint 
- Stateless
- Compares results and metadata and emits as metrics


## Not taken into account**
- authentication
- TLS
- retries


## How to run?
Builds are run through Makefile.

```sh
$ make build
$ ./bin/juggler
```
