use tse;
use sgx_types::*;

pub fn report() -> Result<sgx_report_t, sgx_status_t> {
    let target_info = sgx_target_info_t::default();
    let report_data = sgx_report_data_t::default();

    tse::rsgx_create_report(&target_info, &report_data)
}
