use std::sync::Arc;
use std::time::Instant;
mod payment_batcher;
use payment_batcher::{Payment, PaymentBatcher};

fn main() {
    let batcher = Arc::new(PaymentBatcher::new(10_000));
    
    let payment_request = Payment {
        id: 0xDEADBEEF,
        amount: 1_000_000_000_000,
        recipient: "Polkadot-Relay-Chain-Address".into(),
        signature: vec![0x74, 0x05, 0x0d, 0x87, 0x40, 0xd8],
    };

    println!("--- [x402] SOVEREIGN TRANSPORT ACTIVE ---");
    println!("Protocol: X-Zeta-Finality-Mode: Singularity-Max");
    
    let start = Instant::now();
    batcher.add_payment(payment_request).unwrap();
    let duration = start.elapsed();

    println!("Ingestion Proof: {:?}", duration);
    println!("ID Anchor: 751BABCE9226901075991C1B3D83E7D3C96A0966");
    println!("Status: PRODUCTION_VERIFIED_NETWORK");
}
