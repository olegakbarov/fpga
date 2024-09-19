use std::fs::File;
use std::io::prelude::*;

mod config;
mod fabric;
mod gui;
mod place_and_route;
mod simulation;

use config::load_config;
use place_and_route::place_and_route;
use simulation::{SimulationEngine, TimingModel};

use clap::{Parser, Subcommand};

#[derive(Parser)]
#[command(name = "FPGA Emulator")]
#[command(author = "Your Name")]
#[command(version = "1.0")]
#[command(about = "Emulates FPGA designs", long_about = None)]
struct Cli {
    #[command(subcommand)]
    command: Commands,
}

#[derive(Subcommand)]
enum Commands {
    Simulate {
        #[arg(short, long, value_name = "FILE")]
        config: String,

        #[arg(short, long, value_name = "NUMBER", default_value = "1000")]
        cycles: String,

        #[arg(short, long, value_name = "FILE")]
        output: Option<String>,
    },
    Gui {
        #[arg(short, long, value_name = "FILE")]
        config: String,
    },
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let cli = Cli::parse();

    match &cli.command {
        Commands::Simulate {
            config,
            cycles,
            output,
        } => {
            let cycles: u64 = cycles.parse()?;
            run_simulation(config, cycles, output.as_deref())?;
        }
        Commands::Gui { config } => {
            run_gui_server(config).await?;
        }
    }

    Ok(())
}

fn run_simulation(
    config_path: &str,
    cycles: u64,
    output_path: Option<&str>,
) -> Result<(), Box<dyn std::error::Error>> {
    println!("Loading FPGA configuration from {}...", config_path);
    let config = load_config(config_path)?;

    println!("Performing place and route...");
    // let (placement, routing) = place_and_route(&config, 10, 10); // Assuming a 10x10 grid

    // println!("Setting up simulation engine...");
    // let timing_model = TimingModel::default();
    // let mut engine = SimulationEngine::new(config.clone(), timing_model);

    // // Set some initial inputs (you might want to make this configurable)
    // for input_name in &config.inputs {
    //     engine.set_input(input_name, 0);
    // }

    // println!("Running simulation for {} cycles...", cycles);
    // let result = engine.run(cycles as usize);

    // println!("Simulation completed.");
    // if let Some(path) = output_path {
    //     save_simulation_results(&result, path)?;
    //     println!("Results saved to {}", path);
    // } else {
    //     print_simulation_results(&result);
    // }

    Ok(())
}

async fn run_gui_server(config_path: &str) -> Result<(), Box<dyn std::error::Error>> {
    println!("Loading FPGA configuration from {}...", config_path);
    let config = load_config(config_path)?;

    println!("Starting GUI server...");
    gui::start_server(config).await;

    Ok(())
}

fn save_simulation_results(
    result: &simulation::SimulationResult,
    path: &str,
) -> std::io::Result<()> {
    let mut file = File::create(path)?;
    writeln!(file, "Simulation completed in {} cycles", result.cycles)?;
    for (name, values) in &result.outputs {
        write!(file, "{}: ", name)?;
        for value in values {
            write!(file, "{} ", value)?;
        }
        writeln!(file)?;
    }
    Ok(())
}

fn print_simulation_results(result: &simulation::SimulationResult) {
    println!("Simulation completed in {} cycles", result.cycles);
    for (name, values) in &result.outputs {
        print!("{}: ", name);
        for value in values {
            print!("{} ", value);
        }
        println!();
    }
}
