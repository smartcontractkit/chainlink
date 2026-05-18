#![cfg_attr(not(feature = "std"), no_std)]

pub mod worker_profile;
pub mod offchain_worker;
pub mod ledger;
pub mod network_adapters;

use frame_support::{pallet_prelude::*, dispatch::DispatchResult};
use frame_system::pallet_prelude::*;
use sp_std::vec::Vec;

#[frame_support::pallet]
pub mod pallet {
    use super::*;

    #[pallet::config]
    pub trait Config: frame_system::Config {}

    #[pallet::pallet]
    pub struct Pallet<T>(_);

    #[pallet::storage]
    pub(super) type WorkerProfiles<T: Config> = StorageMap<_, Blake2_128Concat, T::AccountId, worker_profile::WorkerInfo, ValueQuery>;

    #[pallet::storage]
    pub(super) type EpochLedger<T: Config> = StorageMap<_, Blake2_128Concat, u64, Vec<u8>, ValueQuery>;

    #[pallet::call]
    impl<T: Config> Pallet<T> {
        #[pallet::weight(10_000)]
        pub fn register_worker(
            origin: OriginFor<T>,
            wage_per_hour: u128,
            pto_allocated: u32
        ) -> DispatchResult {
            let who = ensure_signed(origin)?;
            worker_profile::register_worker::<T>(who, wage_per_hour, pto_allocated)?;
            Ok(())
        }
    }

    #[pallet::hooks]
    impl<T: Config> Hooks<BlockNumberFor<T>> for Pallet<T> {
        fn offchain_worker(block_number: T::BlockNumber) {
            offchain_worker::run::<T>(block_number);
        }
    }
}
