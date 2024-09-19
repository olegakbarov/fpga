use super::FPGAConfig;
use std::fs::File;
use std::io::Read;
use std::path::Path;

pub fn load_config<P: AsRef<Path>>(path: P) -> Result<FPGAConfig, Box<dyn std::error::Error>> {
    let mut file = File::open(path)?;
    let mut contents = String::new();
    file.read_to_string(&mut contents)?;
    let config: FPGAConfig = serde_json::from_str(&contents)?;
    validate_config(&config)?;
    Ok(config)
}

fn validate_config(config: &FPGAConfig) -> Result<(), String> {
    // To ensure each LUT has a unique identifier
    let lut_ids: std::collections::HashSet<_> = config.luts.iter().map(|lut| lut.id).collect();
    if lut_ids.len() != config.luts.len() {
        return Err("Duplicate LUT IDs found".to_string());
    }

    // To ensure each DFF has a unique identifier
    let dff_ids: std::collections::HashSet<_> = config.dffs.iter().map(|dff| dff.id).collect();
    if dff_ids.len() != config.dffs.len() {
        return Err("Duplicate DFF IDs found".to_string());
    }

    // To ensure each BRAM has a unique identifier
    let bram_ids: std::collections::HashSet<_> = config.brams.iter().map(|bram| bram.id).collect();
    if bram_ids.len() != config.brams.len() {
        return Err("Duplicate BRAM IDs found".to_string());
    }

    // To ensure all connections are valid and refer to existing components
    for conn in &config.connections {
        validate_port(&conn.from, &config)?;
        validate_port(&conn.to, &config)?;
    }

    // To ensure each input has a unique name
    let input_set: std::collections::HashSet<_> = config.inputs.iter().collect();
    if input_set.len() != config.inputs.len() {
        return Err("Duplicate input names found".to_string());
    }

    // To ensure each output has a unique name
    let output_set: std::collections::HashSet<_> = config.outputs.iter().collect();
    if output_set.len() != config.outputs.len() {
        return Err("Duplicate output names found".to_string());
    }

    Ok(())
}

fn validate_port(port: &super::ElementPort, config: &FPGAConfig) -> Result<(), String> {
    match port {
        super::ElementPort::Input { name } => {
            if !config.inputs.contains(name) {
                return Err(format!("Input '{}' not found in config", name));
            }
        }
        super::ElementPort::Output { name } => {
            if !config.outputs.contains(name) {
                return Err(format!("Output '{}' not found in config", name));
            }
        }
        super::ElementPort::LUT { id, port } => {
            if !config.luts.iter().any(|lut| lut.id == *id) {
                return Err(format!("LUT with ID {} not found", id));
            }
            if *port >= 4 {
                return Err(format!("Invalid LUT port {} for LUT {}", port, id));
            }
        }
        super::ElementPort::DFF { id, port } => {
            if !config.dffs.iter().any(|dff| dff.id == *id) {
                return Err(format!("DFF with ID {} not found", id));
            }
            if port != "D" && port != "Q" && port != "CLK" {
                return Err(format!("Invalid DFF port {} for DFF {}", port, id));
            }
        }
        super::ElementPort::BRAM { id, port } => {
            if !config.brams.iter().any(|bram| bram.id == *id) {
                return Err(format!("BRAM with ID {} not found", id));
            }
            if port != "address"
                && port != "data_in"
                && port != "data_out"
                && port != "write_enable"
                && port != "clock"
            {
                return Err(format!("Invalid BRAM port {} for BRAM {}", port, id));
            }
        }
    }
    Ok(())
}
