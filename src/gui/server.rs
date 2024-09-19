use crate::config::FPGAConfig;
use crate::fabric::FPGAFabric;
use serde::{Deserialize, Serialize};
use std::sync::{Arc, Mutex};

#[derive(Clone, Debug)]
struct AppState {
    fabric: Arc<Mutex<FPGAFabric>>,
}

#[derive(Deserialize)]
struct SetInputRequest {
    name: String,
    value: u64,
}

#[derive(Serialize)]
struct GetOutputResponse {
    name: String,
    value: u64,
}

#[derive(Serialize)]
struct FPGAState {
    inputs: Vec<GetOutputResponse>,
    outputs: Vec<GetOutputResponse>,
}

pub async fn start_server(config: FPGAConfig) {
    let fabric = FPGAFabric::from_config(config).expect("Failed to create FPGA fabric");
    let state = AppState {
        fabric: Arc::new(Mutex::new(fabric)),
    };

    println!("{:?}", state);

    // let cors = warp::cors()
    //     .allow_any_origin()
    //     .allow_methods(vec!["GET", "POST", "OPTIONS"])
    //     .allow_headers(vec!["Content-Type"]);

    // let get_state = warp::get()
    //     .and(warp::path("api"))
    //     .and(warp::path("state"))
    //     .and(warp::path::end())
    //     .and(state.clone())
    //     .and_then(handle_get_state);

    // let set_input = warp::post()
    //     .and(warp::path("api"))
    //     .and(warp::path("set_input"))
    //     .and(warp::path::end())
    //     .and(warp::body::json())
    //     .and(state.clone())
    //     .and_then(handle_set_input);

    // let evaluate = warp::post()
    //     .and(warp::path("api"))
    //     .and(warp::path("evaluate"))
    //     .and(warp::path::end())
    //     .and(state)
    //     .and_then(handle_evaluate);

    // let routes = get_state.or(set_input).or(evaluate).with(cors);

    println!("Starting server on http://localhost:3030");
    // warp::serve(routes).run(([127, 0, 0, 1], 3030)).await;
}

async fn handle_evaluate(state: AppState) -> Result<impl warp::Reply, warp::Rejection> {
    let mut fabric = state.fabric.lock().unwrap();
    fabric.evaluate();
    Ok(warp::reply::with_status(
        "FPGA evaluated successfully",
        warp::http::StatusCode::OK,
    ))
}

#[derive(Debug)]
struct SetInputError;
impl warp::reject::Reject for SetInputError {}
