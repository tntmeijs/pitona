use log::{error, info};
use simple_logger::SimpleLogger;

use std::{env, time::Duration};

fn main() {
    SimpleLogger::new().init().unwrap();

    if let Some(port) = env::args().nth(1) {
        run(&port);
    } else {
        error!("Please specify the port the daemon should use to connect to the ECU")
    }
}

fn run(port_name: &str) {
    // The ISO 9141-2 is used on a Daytona 675, which has a baud rate of 10.4 kbit/s
    let port_result = serialport::new(port_name, 10_400)
        .timeout(Duration::from_millis(1_000))
        .open();

    if let Err(error) = port_result {
        let all_available_port_names = serialport::available_ports()
            .into_iter()
            .flatten()
            .collect::<Vec<_>>()
            .iter()
            .map(|p| format!("\"{}\"", p.port_name))
            .collect::<Vec<_>>();

        error!("Failed to open serial port \"{}\": {}", port_name, error);
        info!(
            "Available serial ports: {}",
            all_available_port_names.join(", ")
        );
        return;
    }

    let port = port_result.unwrap();

    info!("Daemon listening on serial port \"{}\"", port_name);
}
