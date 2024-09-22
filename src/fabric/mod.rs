mod bram;
mod dff;
mod lut;
mod wire;

pub use bram::BlockRAM;
pub use dff::DFF;
pub use lut::LUT;
pub use wire::Wire;

use crate::config::{ElementPort, FPGAConfig, FPGAElement};
use std::collections::HashMap;

#[derive(Clone, Debug)]
pub struct FPGAFabric {
    inputs: HashMap<String, usize>,
    outputs: HashMap<String, usize>,
    wires: HashMap<usize, Wire>,
    luts: HashMap<usize, LUT>,
    dffs: HashMap<usize, DFF>,
    brams: HashMap<usize, BlockRAM>,
    bram_connections: HashMap<(usize, String), usize>, // (BRAM id, port) -> wire id
    next_wire_id: usize,
}

impl FPGAFabric {
    pub fn new() -> Self {
        FPGAFabric {
            inputs: HashMap::new(),
            outputs: HashMap::new(),
            wires: HashMap::new(),
            luts: HashMap::new(),
            dffs: HashMap::new(),
            brams: HashMap::new(),
            bram_connections: HashMap::new(),
            next_wire_id: 0,
        }
    }

    pub fn from_config(config: FPGAConfig) -> Result<Self, String> {
        let mut fabric = FPGAFabric::new();

        // Create input wires
        for input_name in &config.inputs {
            let wire_id = fabric.create_wire();
            fabric.inputs.insert(input_name.clone(), wire_id);
        }

        // Create output wires
        for output_name in &config.outputs {
            let wire_id = fabric.create_wire();
            fabric.outputs.insert(output_name.clone(), wire_id);
        }

        // Create LUTs
        for lut_config in &config.luts {
            let lut = LUT::new(lut_config.id, lut_config.truth_table);
            fabric.luts.insert(lut_config.id, lut);
        }

        // Create DFFs
        for dff_config in &config.dffs {
            let dff = DFF::new(dff_config.id);
            fabric.dffs.insert(dff_config.id, dff);
        }

        // Create BRAMs
        for bram_config in &config.brams {
            let bram = BlockRAM::new(bram_config.id, bram_config.size, bram_config.width);
            fabric.brams.insert(bram_config.id, bram);

            // Create wires for BRAM ports
            for port in &["address", "data_in", "data_out", "write_enable", "clock"] {
                let wire_id = fabric.create_wire();
                fabric
                    .bram_connections
                    .insert((bram_config.id, port.to_string()), wire_id);
            }
        }

        // Connect components
        for conn in &config.connections {
            fabric.add_connection(&conn.from, &conn.to)?;
        }

        Ok(fabric)
    }

    fn create_wire(&mut self) -> usize {
        let wire_id = self.next_wire_id;
        self.wires.insert(wire_id, Wire::new(wire_id));
        self.next_wire_id += 1;
        wire_id
    }

    fn add_connection(&mut self, from: &ElementPort, to: &ElementPort) -> Result<(), String> {
        let from_wire = self.get_or_create_wire(from)?;

        match to {
            ElementPort::LUT { id, port } => {
                let lut = self
                    .luts
                    .get_mut(id)
                    .ok_or_else(|| format!("LUT with id {} not found", id))?;
                lut.connect_input(*port, from_wire);
            }
            ElementPort::DFF { id, port } => {
                let dff = self
                    .dffs
                    .get_mut(id)
                    .ok_or_else(|| format!("DFF with id {} not found", id))?;
                match port.as_str() {
                    "D" => dff.connect_input(from_wire),
                    "CLK" => dff.connect_clock(from_wire),
                    _ => return Err(format!("Invalid DFF port: {}", port)),
                }
            }
            ElementPort::BRAM { id, port } => {
                let bram = self
                    .brams
                    .get_mut(id)
                    .ok_or_else(|| format!("BRAM with id {} not found", id))?;
                let bram_wire = self
                    .bram_connections
                    .get(&(*id, port.to_string()))
                    .ok_or_else(|| format!("BRAM port {} not found for BRAM {}", port, id))?;
                bram.connect_wire(port, *bram_wire);
                self.wires
                    .get_mut(bram_wire)
                    .unwrap()
                    .add_destination(from_wire);
            }
            ElementPort::Output { name } => {
                if let Some(output_wire) = self.outputs.get(name) {
                    self.wires
                        .get_mut(output_wire)
                        .unwrap()
                        .add_destination(from_wire);
                } else {
                    return Err(format!("Output {} not found", name));
                }
            }
            ElementPort::Input { .. } => return Err("Cannot connect to an input".to_string()),
        }

        Ok(())
    }

    fn get_or_create_wire(&self, port: &ElementPort) -> Result<usize, String> {
        match port {
            ElementPort::Input { name } => self
                .inputs
                .get(name)
                .cloned()
                .ok_or_else(|| format!("Input {} not found", name)),
            ElementPort::Output { name } => self
                .outputs
                .get(name)
                .cloned()
                .ok_or_else(|| format!("Output {} not found", name)),
            ElementPort::LUT { id, .. } => {
                let lut = self
                    .luts
                    .get(id)
                    .ok_or_else(|| format!("LUT with id {} not found", id))?;
                Ok(lut.output_wire())
            }
            ElementPort::DFF { id, port } => {
                let dff = self
                    .dffs
                    .get(id)
                    .ok_or_else(|| format!("DFF with id {} not found", id))?;
                match port.as_str() {
                    "Q" => Ok(dff.output_wire()),
                    _ => Err(format!("Invalid DFF output port: {}", port)),
                }
            }
            ElementPort::BRAM { id, port } => self
                .bram_connections
                .get(&(*id, port.to_string()))
                .cloned()
                .ok_or_else(|| format!("BRAM port {} not found for BRAM {}", port, id)),
        }
    }

    pub fn set_input(&mut self, name: &str, value: usize) -> Result<(), String> {
        let wire_id = self
            .inputs
            .get(name)
            .ok_or_else(|| format!("Input {} not found", name))?;
        self.wires.get_mut(wire_id).unwrap().set_value(value);
        Ok(())
    }

    pub fn get_output(&self, name: &str) -> Result<usize, String> {
        let wire_id = self
            .outputs
            .get(name)
            .ok_or_else(|| format!("Output {} not found", name))?;
        Ok(self.wires.get(wire_id).unwrap().value())
    }

    pub fn evaluate(&mut self) {
        // Evaluate LUTs
        for lut in self.luts.values() {
            let output = lut.evaluate(&self.wires);
            let output_wire = lut.output_wire();
            self.wires.get_mut(&output_wire).unwrap().set_value(output);
        }

        // Evaluate BRAMs
        for bram in self.brams.values_mut() {
            bram.evaluate(&mut self.wires);
        }

        // Evaluate DFFs
        for dff in self.dffs.values_mut() {
            dff.evaluate(&mut self.wires);
        }

        // Propagate values to outputs
        let mut output_updates = Vec::new();
        for (_, wire_id) in &self.outputs {
            let value = self.wires[wire_id].value();
            let destinations = self.wires[wire_id].destinations().clone();
            output_updates.push((destinations, value));
        }
        for (destinations, value) in output_updates {
            for dest in destinations {
                self.wires.get_mut(&dest).unwrap().set_value(value);
            }
        }
    }

    pub fn get_input_wire_index(&self, input_name: &str) -> Option<usize> {
        self.inputs.get(input_name).cloned()
    }

    pub fn get_output_names(&self) -> Vec<String> {
        self.outputs.keys().cloned().collect()
    }

    pub fn set_wire_value(&mut self, wire_index: usize, value: usize) {
        if let Some(wire) = self.wires.get_mut(&wire_index) {
            wire.set_value(value);
        }
    }

    pub fn get_affected_elements(&self, wire_index: Option<usize>) -> Vec<FPGAElement> {
        let mut affected = Vec::new();

        if let Some(wire_index) = wire_index {
            for (id, lut) in &self.luts {
                if lut.inputs().contains(&wire_index) {
                    affected.push(FPGAElement::LUT(*id));
                }
            }

            for (id, dff) in &self.dffs {
                if dff.input_wire() == Some(wire_index) || dff.clock_wire() == Some(wire_index) {
                    affected.push(FPGAElement::DFF(*id));
                }
            }

            for (id, _) in &self.brams {
                if self.bram_connections.values().any(|&w| w == wire_index) {
                    affected.push(FPGAElement::BRAM(*id));
                }
            }

            for (name, &wire) in &self.outputs {
                if wire == wire_index {
                    affected.push(FPGAElement::Output(name.clone()));
                }
            }
        }

        affected
    }
    pub fn evaluate_lut(&self, lut_id: usize) -> Option<usize> {
        self.luts.get(&lut_id).map(|lut| lut.evaluate(&self.wires))
    }

    pub fn get_lut_output_wire(&self, lut_id: usize) -> usize {
        self.luts
            .get(&lut_id)
            .map(|lut| lut.output_wire())
            .unwrap_or(0)
    }

    pub fn get_bram_wires(&self, bram_id: usize) -> Option<(usize, usize, usize)> {
        let address_wire = self
            .bram_connections
            .get(&(bram_id, "address".to_string()))?;
        let data_in_wire = self
            .bram_connections
            .get(&(bram_id, "data_in".to_string()))?;
        let data_out_wire = self
            .bram_connections
            .get(&(bram_id, "data_out".to_string()))?;
        Some((*address_wire, *data_in_wire, *data_out_wire))
    }

    pub fn get_wire_value(&self, wire_index: usize) -> usize {
        self.wires
            .get(&wire_index)
            .map(|wire| wire.value())
            .unwrap_or(0)
    }

    pub fn write_bram(&mut self, bram_id: usize, address: usize, data: usize) {
        if let Some(bram) = self.brams.get_mut(&bram_id) {
            bram.write(address, data);
        }
    }

    pub fn read_bram(&self, bram_id: usize, address: usize) -> usize {
        self.brams
            .get(&bram_id)
            .map(|bram| bram.read(address))
            .unwrap_or(0)
    }

    pub fn evaluate_dff(&mut self) {
        for dff in self.dffs.values_mut() {
            dff.evaluate(&mut self.wires);
        }
    }

    pub fn get_all_dffs(&self) -> Vec<&DFF> {
        self.dffs.values().collect()
    }

    pub fn get_bram_output_wire(&self, bram_id: usize) -> Option<usize> {
        self.bram_connections
            .get(&(bram_id, "data_out".to_string()))
            .cloned()
    }

    pub fn get_dff_output_wire(&self, dff_id: usize) -> Option<usize> {
        self.dffs.get(&dff_id).map(|dff| dff.output_wire())
    }

    pub fn get_all_inputs(&self) -> Vec<String> {
        self.inputs.keys().cloned().collect()
    }

    pub fn get_dff_output(&self, dff_id: usize) -> Option<usize> {
        self.dffs.get(&dff_id).and_then(|dff| {
            let output_wire = dff.output_wire();
            self.wires.get(&output_wire).map(|wire| wire.value())
        })
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::config::LUTConfig;

    #[test]
    fn test_fabric_creation_and_evaluation() {
        let config = FPGAConfig {
            luts: vec![LUTConfig {
                id: 0,
                truth_table: [
                    false, true, false, true, false, true, false, true, false, true, false, true,
                    false, true, false, true,
                ],
            }],
            dffs: vec![],
            brams: vec![],
            connections: vec![],
            inputs: vec!["in1".to_string(), "in2".to_string()],
            outputs: vec!["out1".to_string()],
        };

        let mut fabric = FPGAFabric::from_config(config).unwrap();

        // Set inputs
        fabric.set_input("in1", 1).unwrap();
        fabric.set_input("in2", 0).unwrap();

        // Evaluate the fabric
        fabric.evaluate();

        // Check output
        let output = fabric.get_output("out1").unwrap();
        assert_eq!(output, 1);
    }
}
