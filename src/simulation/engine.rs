use super::timing::TimingModel;
use super::SimulationResult;
use crate::config::{FPGAConfig, FPGAElement};
use crate::fabric::FPGAFabric;
use std::collections::VecDeque;

pub struct SimulationEngine {
    fabric: FPGAFabric,
    timing_model: TimingModel,
    event_queue: VecDeque<SimulationEvent>,
    current_time: u64,
}

struct SimulationEvent {
    time: u64,
    wire_index: usize,
    new_value: usize,
}

impl SimulationEngine {
    pub fn new(config: FPGAConfig, timing_model: TimingModel) -> Self {
        SimulationEngine {
            fabric: FPGAFabric::from_config(config).expect("Failed to create FPGA fabric"),
            timing_model,
            event_queue: VecDeque::new(),
            current_time: 0,
        }
    }

    pub fn set_input(&mut self, input_name: &str, value: usize) {
        let wire_index = self.fabric.get_input_wire_index(input_name);
        self.schedule_event(self.current_time, wire_index, value);
    }

    pub fn run(&mut self, max_cycles: usize) -> SimulationResult {
        let mut output_history: Vec<(String, Vec<usize>)> = self
            .fabric
            .get_output_names()
            .into_iter()
            .map(|name| (name, Vec::new()))
            .collect();

        let mut cycle_count = 0;

        while cycle_count < max_cycles {
            if self.event_queue.is_empty() {
                // No more events, advance to next clock cycle
                self.current_time += self.timing_model.clock_period;
                cycle_count += 1;
                self.trigger_clock_events();
            } else {
                let event = self.event_queue.pop_front().unwrap();
                self.current_time = event.time;
                self.process_event(event);
            }

            // Record outputs
            for (name, history) in output_history.iter_mut() {
                if let Ok(value) = self.fabric.get_output(name) {
                    history.push(value);
                }
            }
        }

        SimulationResult {
            cycles: cycle_count,
            outputs: output_history,
        }
    }

    fn schedule_event(&mut self, time: u64, wire_index: Option<usize>, new_value: usize) {
        if let Some(wire_index) = wire_index {
            let event = SimulationEvent {
                time,
                wire_index,
                new_value,
            };
            let insert_position = self
                .event_queue
                .binary_search_by_key(&event.time, |e| e.time)
                .unwrap_or_else(|pos| pos);
            self.event_queue.insert(insert_position, event)
        }
    }

    fn process_event(&mut self, event: SimulationEvent) {
        self.fabric
            .set_wire_value(event.wire_index, event.new_value);
        let affected_elements = self.fabric.get_affected_elements(Some(event.wire_index));

        for element in affected_elements {
            match element {
                FPGAElement::LUT(lut_id) => self.evaluate_lut(lut_id),
                FPGAElement::DFF(dff_id) => self.evaluate_dff(dff_id),
                FPGAElement::BRAM(bram_id) => self.evaluate_bram(bram_id),
                FPGAElement::Output(_) => todo!("TODO"),
            }
        }
    }

    fn evaluate_lut(&mut self, lut_id: usize) {
        if let Some(new_value) = self.fabric.evaluate_lut(lut_id) {
            let output_wire = self.fabric.get_lut_output_wire(lut_id);
            let delay = self.timing_model.lut_delay;
            self.schedule_event(self.current_time + delay, Some(output_wire), new_value);
        }
    }

    fn evaluate_dff(&mut self, _dff_id: usize) {
        // DFFs are evaluated on clock events, so we don't need to do anything here
    }

    fn evaluate_bram(&mut self, bram_id: usize) {
        if let Some((address_wire, data_wire, write_enable_wire)) =
            self.fabric.get_bram_wires(bram_id)
        {
            let address = self.fabric.get_wire_value(address_wire);
            let data = self.fabric.get_wire_value(data_wire);
            let write_enable = self.fabric.get_wire_value(write_enable_wire) != 0;

            if write_enable {
                self.fabric.write_bram(bram_id, address as usize, data);
            } else {
                let read_data = self.fabric.read_bram(bram_id, address as usize);
                let output_wire = self.fabric.get_bram_output_wire(bram_id);
                let delay = self.timing_model.bram_read_delay;
                self.schedule_event(self.current_time + delay, output_wire, read_data);
            }
        }
    }

    fn trigger_clock_events(&mut self) {
        // This function is responsible for simulating clock events in the FPGA
        // It iterates through all D flip-flops (DFFs) in the fabric
        let dff_ids: Vec<usize> = self
            .fabric
            .get_all_dffs()
            .iter()
            .map(|dff| dff.id)
            .collect();

        for &dff_id in &dff_ids {
            // Evaluate each DFF to see if its output should change
            self.fabric.evaluate_dff();
            // Get the new value of the DFF after evaluation
            let new_value = self.fabric.get_dff_output(dff_id);
            // Get the output wire for this DFF
            if let Some(output_wire) = self.fabric.get_dff_output_wire(dff_id) {
                // Get the delay from clock to output (clock-to-Q delay)
                let delay = self.timing_model.dff_clock_to_q;
                // Schedule an event to update the DFF's output wire
                // The event will occur after the clock-to-Q delay
                if let Some(new_value) = new_value {
                    self.schedule_event(self.current_time + delay, Some(output_wire), new_value);
                }
            }
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::config::{ConnectionConfig, DFFConfig, ElementPort, LUTConfig};

    #[test]
    fn test_simple_simulation() {
        let config = FPGAConfig {
            luts: vec![LUTConfig {
                id: 0,
                truth_table: [
                    false, true, false, true, false, true, false, true, false, true, false, true,
                    false, true, false, true,
                ],
            }],
            dffs: vec![DFFConfig { id: 0 }],
            brams: vec![],
            connections: vec![
                ConnectionConfig {
                    from: ElementPort::Input {
                        name: "in1".to_string(),
                    },
                    to: ElementPort::LUT { id: 0, port: 0 },
                },
                ConnectionConfig {
                    from: ElementPort::LUT { id: 0, port: 0 },
                    to: ElementPort::DFF {
                        id: 0,
                        port: "D".to_string(),
                    },
                },
                ConnectionConfig {
                    from: ElementPort::DFF {
                        id: 0,
                        port: "Q".to_string(),
                    },
                    to: ElementPort::Output {
                        name: "out1".to_string(),
                    },
                },
            ],
            inputs: vec!["in1".to_string()],
            outputs: vec!["out1".to_string()],
        };

        let timing_model = TimingModel::default();
        let mut engine = SimulationEngine::new(config, timing_model);

        engine.set_input("in1", 1);
        let result = engine.run(10);

        assert_eq!(result.cycles, 10);
        assert_eq!(result.outputs.len(), 1);
        assert_eq!(result.outputs[0].0, "out1");
        assert_eq!(result.outputs[0].1, vec![0, 1, 1, 1, 1, 1, 1, 1, 1, 1]);
    }
}
