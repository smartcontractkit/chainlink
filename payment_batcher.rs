use serde::{Serialize, Deserialize};
use parking_lot::RwLock;
use std::sync::Arc;

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Payment {
    pub id: u64,
    pub amount: u128,
    pub recipient: String,
    pub signature: Vec<u8>,
}

pub struct PaymentBatcher {
    pub queue: Arc<RwLock<Vec<Payment>>>,
    pub capacity: usize,
}

impl PaymentBatcher {
    pub fn new(capacity: usize) -> Self {
        Self {
            queue: Arc::new(RwLock::new(Vec::with_capacity(capacity))),
            capacity,
        }
    }

    pub fn add_payment(&self, payment: Payment) -> Result<(), String> {
        let mut lock = self.queue.write();
        if lock.len() >= self.capacity {
            return Err("Zeta-Singularity pressure release triggered".into());
        }
        lock.push(payment);
        Ok(())
    }
}
