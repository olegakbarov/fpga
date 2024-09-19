use crate::fabric::Wire;
use std::collections::HashMap;

#[derive(Clone, Debug)]
pub struct BlockRAM {
    id: usize,
    size: usize,
    width: usize,
    memory: Vec<usize>,
    address_wire: Option<usize>,
    data_in_wire: Option<usize>,
    data_out_wire: Option<usize>,
    write_enable_wire: Option<usize>,
    clock_wire: Option<usize>,
    last_clock_state: bool,
}

impl BlockRAM {
    pub fn new(id: usize, size: usize, width: usize) -> Self {
        assert!(width <= 64, "BRAM width must be 64 bits or less");
        BlockRAM {
            id,
            size,
            width,
            memory: vec![0; size],
            address_wire: None,
            data_in_wire: None,
            data_out_wire: None,
            write_enable_wire: None,
            clock_wire: None,
            last_clock_state: false,
        }
    }

    pub fn connect_wire(&mut self, port: &str, wire_index: usize) {
        match port {
            "address" => self.address_wire = Some(wire_index),
            "data_in" => self.data_in_wire = Some(wire_index),
            "data_out" => self.data_out_wire = Some(wire_index),
            "write_enable" => self.write_enable_wire = Some(wire_index),
            "clock" => self.clock_wire = Some(wire_index),
            _ => panic!("Invalid BRAM port: {}", port),
        }
    }

    pub fn evaluate(&mut self, wires: &mut HashMap<usize, Wire>) {
        if let Some(clock_wire) = self.clock_wire {
            let clock_state = wires[&clock_wire].value != 0;

            // Detect rising edge
            if clock_state && !self.last_clock_state {
                let address = self.get_address(wires);
                let write_enable = self.get_write_enable(wires);

                if write_enable {
                    let data_in = self.get_data_in(wires);
                    self.write(address, data_in);
                }

                let data_out = self.read(address);
                self.set_data_out(wires, data_out);
            }

            self.last_clock_state = clock_state;
        }
    }

    fn get_address(&self, wires: &HashMap<usize, Wire>) -> usize {
        if let Some(address_wire) = self.address_wire {
            wires[&address_wire].value % self.size
        } else {
            0
        }
    }

    fn get_write_enable(&self, wires: &HashMap<usize, Wire>) -> bool {
        if let Some(write_enable_wire) = self.write_enable_wire {
            wires[&write_enable_wire].value != 0
        } else {
            false
        }
    }

    fn get_data_in(&self, wires: &HashMap<usize, Wire>) -> usize {
        if let Some(data_in_wire) = self.data_in_wire {
            wires[&data_in_wire].value & ((1 << self.width) - 1)
        } else {
            0
        }
    }

    fn set_data_out(&self, wires: &mut HashMap<usize, Wire>, value: usize) {
        if let Some(data_out_wire) = self.data_out_wire {
            wires.get_mut(&data_out_wire).unwrap().value = value;
        }
    }

    pub fn read(&self, address: usize) -> usize {
        self.memory[address] & ((1 << self.width) - 1)
    }

    pub fn write(&mut self, address: usize, data: usize) {
        self.memory[address] = data & ((1 << self.width) - 1);
    }

    // Helper method for testing and debugging
    pub fn get_memory(&self) -> &[usize] {
        &self.memory
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    fn setup_bram_and_wires() -> (BlockRAM, HashMap<usize, Wire>) {
        let mut bram = BlockRAM::new(0, 1024, 32);
        let mut wires = HashMap::new();

        for i in 0..5 {
            wires.insert(i, Wire::new(i));
        }

        bram.connect_wire("address", 0);
        bram.connect_wire("data_in", 1);
        bram.connect_wire("data_out", 2);
        bram.connect_wire("write_enable", 3);
        bram.connect_wire("clock", 4);

        (bram, wires)
    }

    #[test]
    fn test_bram_write_and_read() {
        let (mut bram, mut wires) = setup_bram_and_wires();

        // Set up write operation
        wires.get_mut(&0).unwrap().value = 42; // address
        wires.get_mut(&1).unwrap().value = 0xDEADBEEF; // data_in
        wires.get_mut(&3).unwrap().value = 1; // write_enable

        // Simulate clock edge
        wires.get_mut(&4).unwrap().value = 0;
        bram.evaluate(&mut wires);
        wires.get_mut(&4).unwrap().value = 1;
        bram.evaluate(&mut wires);

        // Verify write
        assert_eq!(bram.get_memory()[42], 0xDEADBEEF);

        // Set up read operation
        wires.get_mut(&3).unwrap().value = 0; // write_enable off

        // Simulate clock edge
        wires.get_mut(&4).unwrap().value = 0;
        bram.evaluate(&mut wires);
        wires.get_mut(&4).unwrap().value = 1;
        bram.evaluate(&mut wires);

        // Verify read
        assert_eq!(wires[&2].value, 0xDEADBEEF);
    }

    #[test]
    fn test_bram_width_limitation() {
        let (mut bram, mut wires) = setup_bram_and_wires();

        // Set up write operation with value exceeding width
        wires.get_mut(&0).unwrap().value = 0; // address
        wires.get_mut(&1).unwrap().value = 0xFFFFFFFFFFFFFFFF; // data_in (64 bits)
        wires.get_mut(&3).unwrap().value = 1; // write_enable

        // Simulate clock edge
        wires.get_mut(&4).unwrap().value = 0;
        bram.evaluate(&mut wires);
        wires.get_mut(&4).unwrap().value = 1;
        bram.evaluate(&mut wires);

        // Verify write (should be truncated to 32 bits)
        assert_eq!(bram.get_memory()[0], 0xFFFFFFFF);
    }
}
