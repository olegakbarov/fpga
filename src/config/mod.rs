mod loader;

pub use loader::load_config;

use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use std::hash::{Hash, Hasher};

#[derive(Deserialize, Serialize, Clone, Debug)]
pub struct FPGAConfig {
    pub luts: Vec<LUTConfig>,
    pub dffs: Vec<DFFConfig>,
    pub brams: Vec<BRAMConfig>,
    pub connections: Vec<ConnectionConfig>,
    pub inputs: Vec<String>,
    pub outputs: Vec<String>,
}

#[derive(Deserialize, Serialize, Clone, Debug)]
pub struct LUTConfig {
    pub id: usize,
    pub truth_table: [bool; 16],
}

#[derive(Deserialize, Serialize, Clone, Debug)]
pub struct DFFConfig {
    pub id: usize,
}

#[derive(Deserialize, Serialize, Clone, Debug)]
pub struct BRAMConfig {
    pub id: usize,
    pub size: usize,
    pub width: usize,
    pub connections: HashMap<String, ElementPort>,
}

#[derive(Deserialize, Serialize, Clone, Debug)]
pub struct ConnectionConfig {
    pub from: ElementPort,
    pub to: ElementPort,
}

#[derive(Deserialize, Serialize, Clone, Debug)]
pub enum ElementPort {
    Input { name: String },
    Output { name: String },
    LUT { id: usize, port: usize },
    DFF { id: usize, port: String },
    BRAM { id: usize, port: String },
}

impl PartialEq for ElementPort {
    fn eq(&self, other: &Self) -> bool {
        match (self, other) {
            (ElementPort::Input { name: n1 }, ElementPort::Input { name: n2 }) => n1 == n2,
            (ElementPort::Output { name: n1 }, ElementPort::Output { name: n2 }) => n1 == n2,
            (ElementPort::LUT { id: id1, port: p1 }, ElementPort::LUT { id: id2, port: p2 }) => {
                id1 == id2 && p1 == p2
            }
            (ElementPort::DFF { id: id1, port: p1 }, ElementPort::DFF { id: id2, port: p2 }) => {
                id1 == id2 && p1 == p2
            }
            (ElementPort::BRAM { id: id1, port: p1 }, ElementPort::BRAM { id: id2, port: p2 }) => {
                id1 == id2 && p1 == p2
            }
            _ => false,
        }
    }
}

impl Eq for ElementPort {}

impl Hash for ElementPort {
    fn hash<H: Hasher>(&self, state: &mut H) {
        match self {
            ElementPort::Input { name } => {
                "Input".hash(state);
                name.hash(state);
            }
            ElementPort::Output { name } => {
                "Output".hash(state);
                name.hash(state);
            }
            ElementPort::LUT { id, port } => {
                "LUT".hash(state);
                id.hash(state);
                port.hash(state);
            }
            ElementPort::DFF { id, port } => {
                "DFF".hash(state);
                id.hash(state);
                port.hash(state);
            }
            ElementPort::BRAM { id, port } => {
                "BRAM".hash(state);
                id.hash(state);
                port.hash(state);
            }
        }
    }
}

pub enum FPGAElement {
    LUT(usize),
    DFF(usize),
    BRAM(usize),
    Output(String),
}
