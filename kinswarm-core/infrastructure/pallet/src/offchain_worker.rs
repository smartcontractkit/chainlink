use crate::ledger::*;
use crate::network_adapters::*;
use frame_system::pallet_prelude::*;
use codec::Encode;
use sp_std::vec::Vec;

pub fn run<T: crate::Config>(_block_number: T::BlockNumber) {
    for (account, info) in crate::pallet::WorkerProfiles::<T>::iter() {
        let net_amount = (info.hours_worked as u128 + info.pto_used as u128) * info.wage_per_hour;
        let payload = sign_payload(&account, net_amount);

        let adapters = NetworkAdapters::default();
        for adapter in adapters.list.iter() {
            adapter.send(payload.clone());
        }

        record_ledger::<T>(account.clone(), net_amount, payload);
    }
}

pub fn sign_payload<T: Encode>(account: &T, amount: u128) -> Vec<u8> {
    let entropy = sp_io::hashing::blake2_256(&account.encode());
    let mood = sp_io::hashing::keccak_256(&account.encode());
    [amount.encode(), entropy.to_vec(), mood.to_vec()].concat()
}
