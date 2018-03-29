# Prismatica IdAM

This is the central IdAM connector for Prismatica

# Running

```bash
$ go install github.com/Project-Prismatica/prismatica-idam/go/prismatica-idam-server
$ ls $GOPATH/bin
prismatica-idam-server
$ 
```

# The API

The service exposes two domains of API endpoints: one for
[Ambassador](http://getambassador.io) external authentication and the other for
interfacing with the IdAM service's use of
[Javascript Web Tokens](https://jwt.io) for use with the microservice ecosystem.

# JWT API

```TODO```

# Ambassador Authentication

Ambassador should use the prefix path ```/ambassador/extauth``` to forward
external requests and allow the header ```x-prismatica-session``. An example
configuration is:

```yaml
---
apiVersion: ambassador/v0
kind:  Module
name:  authentication
config:
  auth_service: "prismatica-idam:8080"
  path_prefix: "/ambassador/extauth/"
  allowed_headers:
  - "x-prismatica-session"
```
