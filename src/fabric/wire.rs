use std::collections::HashMap;

#[derive(Clone, Debug)]
pub struct Wire {
    id: usize,
    pub value: usize,
    new_value: usize,
    destinations: Vec<usize>,
}

impl Wire {
    pub fn new(id: usize) -> Self {
        Wire {
            id,
            value: 0,
            new_value: 0,
            destinations: Vec::new(),
        }
    }

    pub fn id(&self) -> usize {
        self.id
    }

    pub fn value(&self) -> usize {
        self.value
    }

    pub fn set_value(&mut self, value: usize) {
        self.new_value = value;
    }

    pub fn add_destination(&mut self, destination: usize) {
        self.destinations.push(destination);
    }

    pub fn propagate(&mut self, wires: &mut HashMap<usize, Wire>) {
        if self.new_value != self.value {
            self.value = self.new_value;
            for &dest_id in &self.destinations {
                if let Some(dest_wire) = wires.get_mut(&dest_id) {
                    dest_wire.new_value = self.value;
                }
            }
        }
    }

    pub fn destinations(&self) -> &Vec<usize> {
        &self.destinations
    }
}

impl PartialEq for Wire {
    fn eq(&self, other: &Self) -> bool {
        self.id == other.id
    }
}

impl Eq for Wire {}
