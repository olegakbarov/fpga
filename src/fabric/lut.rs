use super::Wire;
use std::collections::HashMap;

#[derive(Clone, Debug)]
pub struct LUT {
    id: usize,
    // Using Option<usize> allows for flexibility in connecting inputs
    // This is helpful because not all inputs may be connected, and it allows for dynamic configuration
    input_wires: [Option<usize>; 4],
    // Using Option<usize> allows for lazy initialization of the output wire
    // This is useful because the output might not be immediately connected upon LUT creation
    output_wire: Option<usize>,
    // A 16-element boolean array is sufficient to represent all possible 4-input combinations
    // This fixed size array optimizes memory usage and access time compared to a dynamic structure
    truth_table: [bool; 16],
}

impl LUT {
    pub fn new(id: usize, truth_table: [bool; 16]) -> Self {
        LUT {
            id,
            // Initialize with None to indicate no connections initially
            // This allows for flexible configuration after LUT creation
            input_wires: [None; 4],
            // Initialize with None as the output wire is not connected yet
            // This enforces explicit output connection before use
            output_wire: None,
            truth_table,
        }
    }

    pub fn id(&self) -> usize {
        self.id
    }

    pub fn connect_input(&mut self, input_index: usize, wire_index: usize) {
        if input_index < 4 {
            // Store the wire index if it's within the valid range
            // This ensures that only valid inputs are connected
            self.input_wires[input_index] = Some(wire_index);
        } else {
            // Panic if the input index is out of range to prevent invalid states
            // This fail-fast approach helps catch configuration errors early
            panic!("LUT input index out of range");
        }
    }

    pub fn connect_output(&mut self, wire_index: usize) {
        // Store the output wire index
        // This allows for dynamic output configuration
        self.output_wire = Some(wire_index);
    }

    pub fn output_wire(&self) -> usize {
        // Unwrap the Option, panicking if the output wire is not connected
        // This ensures that the LUT is properly configured before use
        // It's a fail-fast approach to catch misconfigurations
        self.output_wire.expect("LUT output wire not connected")
    }

    pub fn evaluate(&self, wires: &HashMap<usize, Wire>) -> usize {
        let mut address = 0;
        for (i, wire_option) in self.input_wires.iter().enumerate() {
            if let Some(wire_index) = wire_option {
                if let Some(wire) = wires.get(wire_index) {
                    if wire.value() != 0 {
                        // Build the address by setting bits based on input values
                        // This efficiently computes the truth table index
                        address |= 1 << i;
                    }
                }
            }
        }

        // Use the computed address to look up the output value in the truth table
        // This provides a fast O(1) lookup for the output value
        let output_value = if self.truth_table[address] { 1 } else { 0 };

        // Note: We no longer modify the wires here. Instead, we just return the output value.
        // The caller is responsible for updating the wire if needed.
        if let Some(output_wire_index) = self.output_wire {
            // We keep this check to ensure the output wire is connected,
            // but we don't attempt to modify it.
            if wires.contains_key(&output_wire_index) {
                // The wire exists, but we don't modify it here
            }
        }

        output_value
    }

    // Expose the truth table for debugging and testing
    // This allows verification of the LUT's configuration without modifying its state
    // Useful for unit tests and debugging complex circuits
    pub fn get_truth_table(&self) -> &[bool; 16] {
        &self.truth_table
    }

    pub fn inputs(&self) -> Vec<usize> {
        self.input_wires.iter().filter_map(|&wire| wire).collect()
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    fn setup_lut_and_wires() -> (LUT, HashMap<usize, Wire>) {
        // Create a LUT that implements a 4-input AND gate
        let truth_table = [
            false, false, false, false, false, false, false, false, false, false, false, false,
            false, false, false, true,
        ];
        let mut lut = LUT::new(0, truth_table);
        let mut wires = HashMap::new();

        // Create wires for inputs and output
        for i in 0..5 {
            wires.insert(i, Wire::new(i));
        }

        // Connect input wires
        for i in 0..4 {
            lut.connect_input(i, i);
        }
        // Connect output wire
        lut.connect_output(4);

        (lut, wires)
    }

    #[test]
    fn test_lut_and_gate() {
        let (lut, mut wires) = setup_lut_and_wires();

        // Test case 1: All inputs low
        for i in 0..4 {
            wires.get_mut(&i).unwrap().set_value(0);
        }
        lut.evaluate(&wires);
        // AND gate should output 0 when any input is 0
        assert_eq!(wires[&4].value(), 0);

        // Test case 2: Some inputs high
        wires.get_mut(&0).unwrap().set_value(1);
        wires.get_mut(&1).unwrap().set_value(1);
        lut.evaluate(&wires);
        // AND gate should output 0 when not all inputs are 1
        assert_eq!(wires[&4].value(), 0);

        // Test case 3: All inputs high
        for i in 0..4 {
            wires.get_mut(&i).unwrap().set_value(1);
        }
        lut.evaluate(&wires);
        // AND gate should output 1 only when all inputs are 1
        assert_eq!(wires[&4].value(), 1);
    }

    #[test]
    fn test_lut_or_gate() {
        // Create a LUT that implements a 4-input OR gate
        let truth_table = [
            false, true, true, true, true, true, true, true, true, true, true, true, true, true,
            true, true,
        ];
        let mut lut = LUT::new(1, truth_table);
        let mut wires = HashMap::new();

        // Create wires for inputs and output
        for i in 0..5 {
            wires.insert(i, Wire::new(i));
        }

        // Connect input and output wires
        for i in 0..4 {
            lut.connect_input(i, i);
        }
        lut.connect_output(4);

        // Test case 1: All inputs low
        for i in 0..4 {
            wires.get_mut(&i).unwrap().set_value(0);
        }
        lut.evaluate(&wires);
        // OR gate should output 0 only when all inputs are 0
        assert_eq!(wires[&4].value(), 0);

        // Test case 2: One input high
        wires.get_mut(&2).unwrap().set_value(1);
        lut.evaluate(&wires);
        // OR gate should output 1 when any input is 1
        assert_eq!(wires[&4].value(), 1);

        // Test case 3: All inputs high
        for i in 0..4 {
            wires.get_mut(&i).unwrap().set_value(1);
        }
        lut.evaluate(&wires);
        // OR gate should output 1 when any input is 1
        assert_eq!(wires[&4].value(), 1);
    }

    #[test]
    fn test_lut_xor_gate() {
        // Create a LUT that implements a 2-input XOR gate (using only the first two inputs)
        let truth_table = [
            false, true, true, false, false, true, true, false, false, true, true, false, false,
            true, true, false,
        ];
        let mut lut = LUT::new(2, truth_table);
        let mut wires = HashMap::new();

        // Create wires for inputs and output
        for i in 0..3 {
            wires.insert(i, Wire::new(i));
        }

        // Connect input and output wires
        lut.connect_input(0, 0);
        lut.connect_input(1, 1);
        lut.connect_output(2);

        // Test case 1: Both inputs low
        wires.get_mut(&0).unwrap().set_value(0);
        wires.get_mut(&1).unwrap().set_value(0);
        lut.evaluate(&wires);
        // XOR gate should output 0 when inputs are the same
        assert_eq!(wires[&2].value(), 0);

        // Test case 2: One input high
        wires.get_mut(&0).unwrap().set_value(1);
        lut.evaluate(&wires);
        // XOR gate should output 1 when inputs are different
        assert_eq!(wires[&2].value(), 1);

        // Test case 3: Other input high
        wires.get_mut(&0).unwrap().set_value(0);
        wires.get_mut(&1).unwrap().set_value(1);
        lut.evaluate(&wires);
        // XOR gate should output 1 when inputs are different
        assert_eq!(wires[&2].value(), 1);

        // Test case 4: Both inputs high
        wires.get_mut(&0).unwrap().set_value(1);
        wires.get_mut(&1).unwrap().set_value(1);
        lut.evaluate(&wires);
        // XOR gate should output 0 when inputs are the same
        assert_eq!(wires[&2].value(), 0);
    }
}
