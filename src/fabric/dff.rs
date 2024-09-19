use super::Wire;
use std::collections::HashMap;

#[derive(Clone, Debug)]
pub struct DFF {
    pub id: usize,
    input_wire: Option<usize>,
    clock_wire: Option<usize>,
    output_wire: Option<usize>,
    state: bool,
    last_clock_state: bool,
}

impl DFF {
    pub fn new(id: usize) -> Self {
        DFF {
            id,
            input_wire: None,
            clock_wire: None,
            output_wire: None,
            state: false,
            last_clock_state: false,
        }
    }

    pub fn id(&self) -> usize {
        self.id
    }

    pub fn connect_input(&mut self, wire_index: usize) {
        self.input_wire = Some(wire_index);
    }

    pub fn connect_clock(&mut self, wire_index: usize) {
        self.clock_wire = Some(wire_index);
    }

    pub fn connect_output(&mut self, wire_index: usize) {
        self.output_wire = Some(wire_index);
    }

    pub fn output_wire(&self) -> usize {
        self.output_wire.expect("DFF output wire not connected")
    }

    pub fn get_state(&self) -> bool {
        self.state
    }

    pub fn input_wire(&self) -> Option<usize> {
        self.input_wire
    }

    pub fn clock_wire(&self) -> Option<usize> {
        self.clock_wire
    }

    pub fn evaluate(&mut self, wires: &mut HashMap<usize, Wire>) {
        if let (Some(input_wire), Some(clock_wire)) = (self.input_wire, self.clock_wire) {
            let input_value = wires[&input_wire].value() != 0;
            let clock_state = wires[&clock_wire].value() != 0;

            // Check for rising edge of clock
            if clock_state && !self.last_clock_state {
                // Update state on rising edge
                self.state = input_value;
            }

            self.last_clock_state = clock_state;

            // Update output wire
            if let Some(output_wire) = self.output_wire {
                if let Some(wire) = wires.get_mut(&output_wire) {
                    wire.set_value(if self.state { 1 } else { 0 });
                }
            }
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    fn setup_dff_and_wires() -> (DFF, HashMap<usize, Wire>) {
        let mut dff = DFF::new(0);
        let mut wires = HashMap::new();

        for i in 0..3 {
            wires.insert(i, Wire::new(i));
        }

        dff.connect_input(0);
        dff.connect_clock(1);
        dff.connect_output(2);

        (dff, wires)
    }

    #[test]
    fn test_dff_rising_edge_trigger() {
        let (mut dff, mut wires) = setup_dff_and_wires();

        // Set initial state
        wires.get_mut(&0).unwrap().set_value(1); // Input high
        wires.get_mut(&1).unwrap().set_value(0); // Clock low
        dff.evaluate(&wires);
        assert_eq!(dff.get_state(), false); // State shouldn't change yet

        // Trigger rising edge
        wires.get_mut(&1).unwrap().set_value(1); // Clock high
        dff.evaluate(&wires);
        assert_eq!(dff.get_state(), true); // State should now be high

        // Change input, but no clock edge
        wires.get_mut(&0).unwrap().set_value(0); // Input low
        dff.evaluate(&wires);
        assert_eq!(dff.get_state(), true); // State should remain high

        // Another rising edge
        wires.get_mut(&1).unwrap().set_value(0); // Clock low
        dff.evaluate(&wires);
        wires.get_mut(&1).unwrap().set_value(1); // Clock high
        dff.evaluate(&wires);
        assert_eq!(dff.get_state(), false); // State should now be low

        // Check output wire
        assert_eq!(wires[&2].value(), 0);
    }

    #[test]
    fn test_dff_no_change_on_falling_edge() {
        let (mut dff, mut wires) = setup_dff_and_wires();

        // Set up initial state
        wires.get_mut(&0).unwrap().set_value(1); // Input high
        wires.get_mut(&1).unwrap().set_value(1); // Clock high
        dff.evaluate(&wires);

        // Trigger falling edge
        wires.get_mut(&1).unwrap().set_value(0); // Clock low
        dff.evaluate(&wires);
        assert_eq!(dff.get_state(), false); // State should not change on falling edge

        // Change input and trigger rising edge
        wires.get_mut(&0).unwrap().set_value(1); // Input high
        wires.get_mut(&1).unwrap().set_value(1); // Clock high
        dff.evaluate(&wires);
        assert_eq!(dff.get_state(), true); // State should change on rising edge
    }
}
