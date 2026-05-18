use codec::{Encode, Decode};
use frame_support::pallet_prelude::TypeInfo;

#[derive(Encode, Decode, Default, Clone, PartialEq, TypeInfo)]
pub struct WorkerInfo {
    pub wage_per_hour: u128,
    pub hours_worked: u32,
    pub pto_allocated: u32,
    pub pto_used: u32,
    pub last_epoch_paid: u64,
}

pub fn register_worker<T: frame_system::Config>(
    account: T::AccountId,
    wage_per_hour: u128,
    pto_allocated: u32
) -> frame_support::dispatch::DispatchResult {
    crate::pallet::WorkerProfiles::<T>::insert(account, WorkerInfo {
        wage_per_hour,
        hours_worked: 0,
        pto_allocated,
        pto_used: 0,
        last_epoch_paid: 0,
    });
    Ok(())
}
