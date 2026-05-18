use sp_std::vec::Vec;
use sp_std::boxed::Box;

pub trait NetworkAdapter {
    fn send(&self, payload: Vec<u8>) -> bool;
}

pub struct NetworkAdapters {
    pub list: Vec<Box<dyn NetworkAdapter>>,
}

impl Default for NetworkAdapters {
    fn default() -> Self {
        Self { list: sp_std::vec![Box::new(EVMAdapter), Box::new(CosmosAdapter), Box::new(SolanaAdapter)] }
    }
}

pub struct EVMAdapter;
impl NetworkAdapter for EVMAdapter {
    fn send(&self, _payload: Vec<u8>) -> bool { true }
}

pub struct CosmosAdapter;
impl NetworkAdapter for CosmosAdapter {
    fn send(&self, _payload: Vec<u8>) -> bool { true }
}

pub struct SolanaAdapter;
impl NetworkAdapter for SolanaAdapter {
    fn send(&self, _payload: Vec<u8>) -> bool { true }
}
