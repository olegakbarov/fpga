mod placer;
mod router;

use crate::config::FPGAConfig;
use placer::PlacementResult;
use router::RoutingResult;

pub struct PlaceAndRouteResult {
    pub placement: PlacementResult,
    pub routing: RoutingResult,
}

pub fn place_and_route(
    config: &FPGAConfig,
    grid_width: usize,
    grid_height: usize,
) -> PlaceAndRouteResult {
    let placement = placer::place_elements(config, grid_width, grid_height);
    let routing = router::route_connections(config, &placement, grid_width, grid_height);
    PlaceAndRouteResult { placement, routing }
}
