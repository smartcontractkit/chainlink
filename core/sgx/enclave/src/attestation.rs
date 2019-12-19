use sgx_rand::{Rng, os::SgxRng};
use sgx_tcrypto as tcrypto;
use sgx_types::*;
use std::io::{self};
use std::ptr;
use std::vec::Vec;
use tse;

extern "C" {
    pub fn ocall_sgx_init_quote ( ret_val : *mut sgx_status_t,
                  ret_ti  : *mut sgx_target_info_t,
                  ret_gid : *mut sgx_epid_group_id_t) -> sgx_status_t;
    pub fn ocall_get_quote (ret_val            : *mut sgx_status_t,
                p_sigrl            : *const u8,
                sigrl_len          : u32,
                p_report           : *const sgx_report_t,
                quote_type         : sgx_quote_sign_type_t,
                p_spid             : *const sgx_spid_t,
                p_nonce            : *const sgx_quote_nonce_t,
                p_qe_report        : *mut sgx_report_t,
                p_quote            : *mut u8,
                maxlen             : u32,
                p_quote_len        : *mut u32) -> sgx_status_t;
}

const RET_QUOTE_BUF_LEN : u32 = 2048;
type QuoteBuf = [u8; RET_QUOTE_BUF_LEN as usize];

pub fn report() -> Result<sgx_report_t, sgx_status_t> {
    let report_data = sgx_report_data_t::default();
    let (target_info, _) = init_quote()?;

    let report = tse::rsgx_create_report(&target_info, &report_data)?;

    let quote_nonce = match quote_nonce() {
        Ok(n) => n,
        Err(_) => return Err(sgx_status_t::SGX_ERROR_UNEXPECTED)
    };
    let (qe_report, quote_buf, buf_len) = quote(&report, &quote_nonce)?;

    match tse::rsgx_verify_report(&qe_report) {
        Ok(()) => println!("rsgx_verify_report passed!"),
        Err(x) => {
            println!("rsgx_verify_report failed with {:?}", x);
            return Err(x);
        },
    }

    if !report_matches_current_platform(&target_info, &qe_report) {
        println!("qe_report does not match current target_info!");
        return Err(sgx_status_t::SGX_ERROR_UNEXPECTED);
    }

    if report_tampered(&qe_report, &quote_buf, buf_len, &quote_nonce) {
        println!("report has been tampered with");
        return Err(sgx_status_t::SGX_ERROR_UNEXPECTED);
    }

    Ok(qe_report)
}

fn init_quote() -> Result<(sgx_target_info_t, sgx_epid_group_id_t), sgx_status_t> {
    let mut target_info = sgx_target_info_t::default();
    let mut epid_group_id = sgx_epid_group_id_t::default();
    let mut rt = sgx_status_t::SGX_ERROR_UNEXPECTED;

    let result = unsafe {
        ocall_sgx_init_quote(&mut rt as *mut sgx_status_t,
                             &mut target_info as *mut sgx_target_info_t,
                             &mut epid_group_id as *mut sgx_epid_group_id_t)
    };

    println!("epid_group_id = {:?}", epid_group_id);

    if result != sgx_status_t::SGX_SUCCESS {
        println!("ocall_sgx_init_quote failed {}", result);
        return Err(result);
    }

    if rt != sgx_status_t::SGX_SUCCESS {
        println!("ocall_sgx_init_quote returned {}", rt);
        return Err(rt);
    }

    return Ok((target_info, epid_group_id))
}

fn quote_nonce() -> io::Result<sgx_quote_nonce_t> {
    let mut quote_nonce = sgx_quote_nonce_t { rand : [0;16] };
    let mut os_rng = SgxRng::new()?;
    os_rng.fill_bytes(&mut quote_nonce.rand);
    Ok(quote_nonce)
}

fn quote(report: &sgx_report_t, quote_nonce: &sgx_quote_nonce_t)
    -> Result<(sgx_report_t, QuoteBuf, u32), sgx_status_t> {
    let mut qe_report = sgx_report_t::default();
    let mut quote_buf = [0; RET_QUOTE_BUF_LEN as usize];
    let mut quote_len : u32 = 0;

    let spid = sgx_spid_t::default();
    let mut rt = sgx_status_t::SGX_ERROR_UNEXPECTED;

    let result = unsafe {
        ocall_get_quote(
            &mut rt as *mut sgx_status_t,
            ptr::null(), // Not using a signature revocation list yet
            0,
            report as * const sgx_report_t,
            sgx_quote_sign_type_t::SGX_LINKABLE_SIGNATURE,
            &spid as *const sgx_spid_t,
            quote_nonce as * const sgx_quote_nonce_t,
            &mut qe_report as *mut sgx_report_t,
            quote_buf.as_mut_ptr(),
            RET_QUOTE_BUF_LEN as u32,
            &mut quote_len as *mut u32)
    };

    if result != sgx_status_t::SGX_SUCCESS {
        println!("ocall_get_quote failed {}", result);
        return Err(result);
    }

    if rt != sgx_status_t::SGX_SUCCESS {
        println!("ocall_get_quote returned {}", rt);
        return Err(rt);
    }

    Ok((qe_report, quote_buf, quote_len))
}

fn report_matches_current_platform(target_info: &sgx_target_info_t, qe_report: &sgx_report_t) -> bool {
    target_info.mr_enclave.m == qe_report.body.mr_enclave.m &&
        target_info.attributes.flags == qe_report.body.attributes.flags &&
        target_info.attributes.xfrm  == qe_report.body.attributes.xfrm
}

fn report_tampered(qe_report: &sgx_report_t, quote_buf: &QuoteBuf, quote_len: u32, quote_nonce: &sgx_quote_nonce_t) -> bool {
    let mut rhs_vec : Vec<u8> = quote_nonce.rand.to_vec();
    rhs_vec.extend(&quote_buf[..quote_len as usize]);
    let rhs_hash = tcrypto::rsgx_sha256_slice(&rhs_vec[..]).unwrap();
    let lhs_hash = &qe_report.body.report_data.d[..32];

    rhs_hash != lhs_hash
}
