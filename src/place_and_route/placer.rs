use crate::config::FPGAConfig;
use std::collections::HashMap;

#[derive(Debug, Clone)]
pub struct PlacementResult {
    pub lut_positions: HashMap<usize, (usize, usize)>,
    pub dff_positions: HashMap<usize, (usize, usize)>,
    pub bram_positions: HashMap<usize, (usize, usize)>,
}

struct FPGAGrid {
    width: usize,
    height: usize,
    cells: Vec<Vec<CellType>>,
}

#[derive(Clone, PartialEq)]
enum CellType {
    Empty,
    LUT(usize),
    DFF(usize),
    BRAM(usize),
}

impl FPGAGrid {
    fn new(width: usize, height: usize) -> Self {
        FPGAGrid {
            width,
            height,
            cells: vec![vec![CellType::Empty; width]; height],
        }
    }

    fn place_element(&mut self, element: CellType) -> Option<(usize, usize)> {
        for y in 0..self.height {
            for x in 0..self.width {
                if let CellType::Empty = self.cells[y][x] {
                    self.cells[y][x] = element.clone();
                    return Some((x, y));
                }
            }
        }
        None
    }
}

pub fn place_elements(
    config: &FPGAConfig,
    grid_width: usize,
    grid_height: usize,
) -> PlacementResult {
    let mut grid = FPGAGrid::new(grid_width, grid_height);
    let mut result = PlacementResult {
        lut_positions: HashMap::new(),
        dff_positions: HashMap::new(),
        bram_positions: HashMap::new(),
    };

    // Place LUTs
    for lut in &config.luts {
        if let Some(pos) = grid.place_element(CellType::LUT(lut.id)) {
            result.lut_positions.insert(lut.id, pos);
        } else {
            panic!("Not enough space to place all LUTs");
        }
    }

    // Place DFFs
    for dff in &config.dffs {
        if let Some(pos) = grid.place_element(CellType::DFF(dff.id)) {
            result.dff_positions.insert(dff.id, pos);
        } else {
            panic!("Not enough space to place all DFFs");
        }
    }

    // Place BRAMs
    for bram in &config.brams {
        if let Some(pos) = grid.place_element(CellType::BRAM(bram.id)) {
            result.bram_positions.insert(bram.id, pos);
        } else {
            panic!("Not enough space to place all BRAMs");
        }
    }

    result
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::config::{BRAMConfig, DFFConfig, LUTConfig};

    #[test]
    fn test_simple_placement() {
        let config = FPGAConfig {
            luts: vec![LUTConfig {
                id: 0,
                truth_table: [false; 16],
            }],
            dffs: vec![DFFConfig { id: 0 }],
            brams: vec![BRAMConfig {
                id: 0,
                size: 1024,
                width: 8,
                connections: HashMap::new(),
            }],
            connections: vec![],
            inputs: vec![],
            outputs: vec![],
        };

        let placement = place_elements(&config, 2, 2);

        assert_eq!(placement.lut_positions.len(), 1);
        assert_eq!(placement.dff_positions.len(), 1);
        assert_eq!(placement.bram_positions.len(), 1);

        assert!(placement.lut_positions.contains_key(&0));
        assert!(placement.dff_positions.contains_key(&0));
        assert!(placement.bram_positions.contains_key(&0));
    }

    #[test]
    #[should_panic(expected = "Not enough space to place all LUTs")]
    fn test_placement_overflow() {
        let config = FPGAConfig {
            luts: vec![
                LUTConfig {
                    id: 0,
                    truth_table: [false; 16],
                },
                LUTConfig {
                    id: 1,
                    truth_table: [false; 16],
                },
            ],
            dffs: vec![],
            brams: vec![],
            connections: vec![],
            inputs: vec![],
            outputs: vec![],
        };

        place_elements(&config, 1, 1); // This should panic
    }
}
