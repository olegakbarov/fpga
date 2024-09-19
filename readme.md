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
