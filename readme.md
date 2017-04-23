### io.confs.api

Dockerized Go application with Postgres database, proxied by nginx

### Prerequisites

- Go
- Docker

### Start

...

This results running server on port `9999`


### Linking containers caveats

Creates environment variables in proxy container, with ip and port info for go container, also creates entries in /etc/hosts with ip info [other container]:[alias in this container]

```yaml
links:
  - confsio:app
```
### Develop

Recompile Go app `run build.sh`

Build Docker image `docker build -t confsio_img .`

Run Docker image with port-forwarding: `docker run -it -p 8080:8080 confsio_img`

Inspect container's ENV variables: `docker inspect -f "{{ .Config.Env }}" container-id`

Copy file from container to host: `docker cp <containerId>:/file/path/within/container /host/path/target`
