use super::placer::PlacementResult;
use crate::config::{ElementPort, FPGAConfig};
use std::collections::{HashMap, VecDeque};

#[derive(Debug, Clone)]
pub struct RoutingResult {
    pub routes: HashMap<(ElementPort, ElementPort), Vec<(usize, usize)>>,
}

struct RoutingGrid {
    width: usize,
    height: usize,
    cells: Vec<Vec<bool>>,
}

impl RoutingGrid {
    fn new(width: usize, height: usize) -> Self {
        RoutingGrid {
            width,
            height,
            cells: vec![vec![false; width]; height],
        }
    }

    fn is_occupied(&self, x: usize, y: usize) -> bool {
        self.cells[y][x]
    }

    fn occupy(&mut self, x: usize, y: usize) {
        self.cells[y][x] = true;
    }
}

pub fn route_connections(
    config: &FPGAConfig,
    placement: &PlacementResult,
    grid_width: usize,
    grid_height: usize,
) -> RoutingResult {
    let mut routing_grid = RoutingGrid::new(grid_width, grid_height);
    let mut result = RoutingResult {
        routes: HashMap::new(),
    };

    for conn in &config.connections {
        let start = get_element_position(&conn.from, placement);
        let end = get_element_position(&conn.to, placement);

        if let (Some(start), Some(end)) = (start, end) {
            if let Some(path) = find_path(start, end, &mut routing_grid) {
                result
                    .routes
                    .insert((conn.from.clone(), conn.to.clone()), path);
            } else {
                panic!(
                    "Unable to route connection from {:?} to {:?}",
                    conn.from, conn.to
                );
            }
        } else {
            panic!("Unable to find position for connection endpoints");
        }
    }

    result
}

fn get_element_position(port: &ElementPort, placement: &PlacementResult) -> Option<(usize, usize)> {
    match port {
        ElementPort::LUT { id, .. } => placement.lut_positions.get(id).cloned(),
        ElementPort::DFF { id, .. } => placement.dff_positions.get(id).cloned(),
        ElementPort::BRAM { id, .. } => placement.bram_positions.get(id).cloned(),
        _ => None,
    }
}

fn find_path(
    start: (usize, usize),
    end: (usize, usize),
    grid: &mut RoutingGrid,
) -> Option<Vec<(usize, usize)>> {
    let mut queue = VecDeque::new();
    let mut visited = HashMap::new();
    let mut came_from = HashMap::new();

    queue.push_back(start);
    visited.insert(start, true);

    while let Some(current) = queue.pop_front() {
        if current == end {
            return Some(reconstruct_path(start, end, came_from));
        }

        for neighbor in get_neighbors(current, grid.width, grid.height) {
            if !visited.contains_key(&neighbor) && !grid.is_occupied(neighbor.0, neighbor.1) {
                queue.push_back(neighbor);
                visited.insert(neighbor, true);
                came_from.insert(neighbor, current);
            }
        }
    }

    None
}

fn get_neighbors(pos: (usize, usize), width: usize, height: usize) -> Vec<(usize, usize)> {
    let (x, y) = pos;
    let mut neighbors = Vec::new();

    if x > 0 {
        neighbors.push((x - 1, y));
    }
    if x < width - 1 {
        neighbors.push((x + 1, y));
    }
    if y > 0 {
        neighbors.push((x, y - 1));
    }
    if y < height - 1 {
        neighbors.push((x, y + 1));
    }

    neighbors
}

fn reconstruct_path(
    start: (usize, usize),
    end: (usize, usize),
    came_from: HashMap<(usize, usize), (usize, usize)>,
) -> Vec<(usize, usize)> {
    let mut path = vec![end];
    let mut current = end;

    while current != start {
        current = *came_from.get(&current).unwrap();
        path.push(current);
    }

    path.reverse();
    path
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::config::{ConnectionConfig, LUTConfig};

    #[test]
    fn test_simple_routing() {
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
            connections: vec![ConnectionConfig {
                from: ElementPort::LUT { id: 0, port: 0 },
                to: ElementPort::LUT { id: 1, port: 0 },
            }],
            inputs: vec![],
            outputs: vec![],
        };

        let placement = PlacementResult {
            lut_positions: vec![(0, (0, 0)), (1, (2, 2))].into_iter().collect(),
            dff_positions: HashMap::new(),
            bram_positions: HashMap::new(),
        };

        let routing = route_connections(&config, &placement, 3, 3);

        assert_eq!(routing.routes.len(), 1);
        assert!(routing.routes.contains_key(&(
            ElementPort::LUT { id: 0, port: 0 },
            ElementPort::LUT { id: 1, port: 0 }
        )));
    }

    #[test]
    #[should_panic(expected = "Unable to route connection")]
    fn test_impossible_routing() {
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
            connections: vec![ConnectionConfig {
                from: ElementPort::LUT { id: 0, port: 0 },
                to: ElementPort::LUT { id: 1, port: 0 },
            }],
            inputs: vec![],
            outputs: vec![],
        };

        let placement = PlacementResult {
            lut_positions: vec![(0, (0, 0)), (1, (1, 1))].into_iter().collect(),
            dff_positions: HashMap::new(),
            bram_positions: HashMap::new(),
        };

        // This should panic because the grid is too small to route the connection
        route_connections(&config, &placement, 2, 2);
    }
}
