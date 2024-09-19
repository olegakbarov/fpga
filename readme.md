<<<<<<< HEAD
# FPGA Emulator

## Overview

This FPGA Emulator is a sophisticated software tool designed to simulate the behavior of Field-Programmable Gate Arrays (FPGAs). It provides a flexible and extensible platform for designing, testing, and analyzing digital circuits without the need for physical hardware.

## Features

- **Configurable FPGA Architecture**: Supports Look-Up Tables (LUTs), D Flip-Flops (DFFs), and Block RAMs (BRAMs).
- **Place and Route**: Implements basic placement and routing algorithms for mapping logical elements to a grid.
- **Simulation Engine**: Offers cycle-accurate simulation with customizable timing models.
- **GUI Interface**: Provides a web-based graphical interface for interactive circuit manipulation and visualization.
- **Command-Line Interface**: Supports batch simulations and integration into automated workflows.
- **Extensible Design**: Modular architecture allows for easy addition of new FPGA elements and features.

## Getting Started

### Prerequisites

- Rust programming language (latest stable version)
- Cargo package manager

### Installation

1. Clone the repository:
   ```
   git clone https://github.com/olegakbarov/fpga.git
   cd fpga
   ```

2. Build the project:
   ```
   cargo build --release
   ```

### Usage

#### Command-Line Interface

Run a simulation:
```
cargo run -- simulate -c path/to/config.json -n 1000 -o output.txt
```

Start the GUI server:
```
cargo run -- gui -c path/to/config.json
```

#### GUI Interface

After starting the GUI server, open a web browser and navigate to `http://localhost:3030` to access the graphical interface.

## Project Structure

- `src/config/`: Configuration parsing and validation
- `src/fabric/`: FPGA fabric implementation (LUTs, DFFs, BRAMs)
- `src/place_and_route/`: Placement and routing algorithms
- `src/simulation/`: Simulation engine and timing models
- `src/gui/`: Web-based graphical user interface
- `src/main.rs`: Command-line interface and entry point

## Configuration Format

The FPGA configuration is specified in JSON format. Here's a simple example:

```json
{
  "luts": [
    {"id": 0, "truth_table": [false, true, false, true, false, true, false, true, false, true, false, true, false, true, false, true]}
  ],
  "dffs": [
    {"id": 0}
  ],
  "brams": [
    {"id": 0, "size": 1024, "width": 8}
  ],
  "connections": [
    {"from": {"Input": {"name": "in1"}}, "to": {"LUT": {"id": 0, "port": 0}}},
    {"from": {"LUT": {"id": 0, "port": 0}}, "to": {"DFF": {"id": 0, "port": "D"}}}
  ],
  "inputs": ["in1"],
  "outputs": ["out1"]
}
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
=======
### go-web-backend

Dockerized API service for web app

### Prerequisites

- Go
- Docker

### Start

### Develop

TODO

ORM: https://upper.io/db.v2/examples

### Ports cheatsheet

Inside containers:

```
1337 — frontend
3000 — grafana
6666 — postgres
8086, 8083 — influxdb
8080 — cadvisor
9999 — api
```

Exposed:

### Linking containers caveats

```yaml
links:
  - api:api
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

### Docker Monitoring

https://www.brianchristner.io/how-to-setup-docker-monitoring/

###  $ Curls

#### edit :id & :user-id

```
curl -H "Content-Type: application/json" -X PUT -d '{ "name": "UPDATED CONF", "start_date": "2017-10-19T08:00:00Z", "end_date": "2017-10-22T08:00:00Z", "description": "yolo", "picture": null, "country": "USA", "city": "SF", "address": "Rodeo drive 1", "category": "big data", "tickets_available": false, "discount_program": false, "min_price": 0, "max_price": 100, "facebook": "", "youtube": "", "twitter": "", "details": {}, "id": :id, "added_by": :user-id}' http://localhost:9999/api/v1/conf/:id
```

#### create

```
curl -H "Content-Type: application/json" -X POST -d '{ "name": "CREATED!CONF", "start_date": "2017-10-19T08:00:00Z", "end_date": "2017-10-22T08:00:00Z", "description": "yolo", "picture": null, "country": "USA", "city": "SF", "address": "Rodeo drive 1", "category": "big data", "tickets_available": false, "discount_program": false, "min_price": 0, "max_price": 100, "facebook": null, "youtube": null, "twitter": null, "details": {}}' http://localhost:9999/api/v1/conf
```
>>>>>>> 40c30a0b5805be05402d06a9a2b129e3a2bb914f
