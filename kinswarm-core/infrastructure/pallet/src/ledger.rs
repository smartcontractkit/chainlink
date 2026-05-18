use sp_std::vec::Vec;
use frame_support::pallet_prelude::*;

pub fn record_ledger<T: crate::Config>(_account: T::AccountId, _amount: u128, payload: Vec<u8>) {
    let mut ledger = crate::pallet::EpochLedger::<T>::get(0u64);
    ledger.extend(payload);
    crate::pallet::EpochLedger::<T>::insert(0u64, ledger);
}
