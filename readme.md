### What

Dockerized Go application with Postgres database, proxied by nginx

### Prerequisites

- Go
- Docker

### Start

```
$ docker-compose build
$ docker-compose up
```
### Develop

Recompile Go app `run build.sh`

Build Docker image `docker build -t confsio_img .`

Run Docker image with port-forwarding: `docker run -it -p 8080:8080 confsio_img`

Inspect container's ENV variables: `docker inspect -f "{{ .Config.Env }}" container-id`

Copy file from container to host:

```
docker cp <containerId>:/file/path/within/container /host/path/target
```
