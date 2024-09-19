mod engine;
mod timing;

pub use engine::SimulationEngine;
pub use timing::TimingModel;

pub struct SimulationResult {
    pub cycles: usize,
    pub outputs: Vec<(String, Vec<usize>)>,
}
