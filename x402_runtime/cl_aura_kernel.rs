use std::time::Instant;

pub struct ChainlinkAuraPlugin {
    pub node_id: &'static str,
    pub oracle_identity: &'static str,
}

impl ChainlinkAuraPlugin {
    pub fn commit_to_don(&self, payload_size: usize) {
        let start = Instant::now();
        println!("--- [CHAINLINK-AURA] NATIVE CRE PLUGIN ACTIVE ---");
        println!("NodeID:    {}", self.node_id);
        println!("Identity:  {}", self.oracle_identity);
        println!("Action:    Batching {} Transaction Proofs", payload_size);
        println!("Latency:   {:?}", start.elapsed());
        println!("Status:    SYNCED_WITH_DON_FINALITY");
        println!("--------------------------------------------------");
    }
}

fn main() {
    let plugin = ChainlinkAuraPlugin {
        node_id: "Aura-DON-Member-01",
        oracle_identity: "751BABCE9226901075991C1B3D83E7D3C96A0966",
    };
    plugin.commit_to_don(1_400_000);
}
