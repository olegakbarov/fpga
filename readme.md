### io.confs.api

Dockerized API service

### Prerequisites

- Go
- Docker

### Start

...

This results running nginx-server on port `9999` proxying requests to App


### Develop

TODO

### Linking containers caveats

```yaml
links:
  - confsio:app
```

This creates environment variables in proxy container, with ip and port info for go container, also creates entries in /etc/hosts with ip info [other container]:[alias in this container]

```yaml
volumes:
    - ./nginx.conf:/etc/nginx/nginx.conf:ro
```

- Conntect host's `./nginx.conf` with container's `nginx.conf`

- `:ro` means read only perms in container

### Docker tips

Run Docker image with port-forwarding: `docker run -it -p 8080:8080 confsio_img`

Inspect container's ENV variables: `docker inspect -f "{{ .Config.Env }}" container-id`

Copy file from container to host: `docker cp <containerId>:/file/path/within/container /host/path/target`

### 'Tests'

edit :id & :user-id

```
curl -H "Content-Type: application/json" -X PUT -d '{ "name": "UPDATED CONF", "start_date": "2017-10-19T08:00:00Z", "end_date": "2017-10-22T08:00:00Z", "description": "yolo", "picture": null, "country": "USA", "city": "SF", "address": "Rodeo drive 1", "category": "big data", "tickets_available": false, "discount_program": false, "min_price": 0, "max_price": 100, "facebook": "", "youtube": "", "twitter": "", "details": {}, "id": :id, "added_by": :user-id}' http://localhost:9999/api/v1/conf/:id```
