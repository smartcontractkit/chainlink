use sgx_rand::{Rng, os::SgxRng};
use sgx_tcrypto::{self as tcrypto, SgxEccHandle};
use sgx_types::*;
use std::io::{self, Read};
use std::ptr;
use std::string::String;
use std::untrusted::fs;
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

fn quote() -> Result<(sgx_target_info_t, sgx_epid_group_id_t), sgx_status_t> {
    let mut target_info : sgx_target_info_t = sgx_target_info_t::default();
    let mut epid_group_id : sgx_epid_group_id_t = sgx_epid_group_id_t::default();
    let mut rt : sgx_status_t = sgx_status_t::SGX_ERROR_UNEXPECTED;

    let result = unsafe {
        ocall_sgx_init_quote(&mut rt as *mut sgx_status_t,
                             &mut target_info as *mut sgx_target_info_t,
                             &mut epid_group_id as *mut sgx_epid_group_id_t)
    };

    println!("epid_group_id = {:?}", epid_group_id);

    if result != sgx_status_t::SGX_SUCCESS {
        return Err(result);
    }

    if rt != sgx_status_t::SGX_SUCCESS {
        return Err(rt);
    }

    return Ok((target_info, epid_group_id))
}

fn new_report(public_key: &sgx_ec256_public_t) -> sgx_report_data_t {
    let mut report_data: sgx_report_data_t = sgx_report_data_t::default();

    // Fill ecc256 public key into report_data
    let mut pub_k_gx = public_key.gx.clone();
    pub_k_gx.reverse();
    let mut pub_k_gy = public_key.gy.clone();
    pub_k_gy.reverse();
    report_data.d[..32].clone_from_slice(&pub_k_gx);
    report_data.d[32..].clone_from_slice(&pub_k_gy);

    report_data
}

fn keypair() -> Result<(sgx_ec256_private_t, sgx_ec256_public_t), sgx_status_t> {
    let ecc_handle = SgxEccHandle::new();
    let _result = ecc_handle.open();
    ecc_handle.create_key_pair()
}

fn quote_nonce() -> io::Result<sgx_quote_nonce_t> {
    let mut quote_nonce = sgx_quote_nonce_t { rand : [0;16] };
    let mut os_rng = SgxRng::new()?;
    os_rng.fill_bytes(&mut quote_nonce.rand);
    Ok(quote_nonce)
}

pub fn report() -> Result<sgx_report_t, sgx_status_t> {
    let (_, public_key) = keypair()?;
    let report_data = new_report(&public_key);
    let (target_info, _) = quote()?;
    let report = tse::rsgx_create_report(&target_info, &report_data)?;

    let quote_nonce = match quote_nonce() {
        Ok(n) => n,
        Err(_) => return Err(sgx_status_t::SGX_ERROR_UNEXPECTED)
    };

    let mut qe_report = sgx_report_t::default();
    const RET_QUOTE_BUF_LEN : u32 = 2048;
    let mut return_quote_buf : [u8; RET_QUOTE_BUF_LEN as usize] = [0;RET_QUOTE_BUF_LEN as usize];
    let mut quote_len : u32 = 0;

    let p_report = &report as * const sgx_report_t;
    let quote_type = sgx_quote_sign_type_t::SGX_LINKABLE_SIGNATURE;
    let spid : sgx_spid_t = sgx_spid_t::default();


    let p_spid = &spid as *const sgx_spid_t;
    let p_nonce = &quote_nonce as * const sgx_quote_nonce_t;
    let p_qe_report = &mut qe_report as *mut sgx_report_t;
    let p_quote = return_quote_buf.as_mut_ptr();
    let maxlen = RET_QUOTE_BUF_LEN;
    let p_quote_len = &mut quote_len as *mut u32;

    let mut rt : sgx_status_t = sgx_status_t::SGX_ERROR_UNEXPECTED;
    let result = unsafe {
        ocall_get_quote(&mut rt as *mut sgx_status_t,
                ptr::null(),
                0,
                p_report,
                quote_type,
                p_spid,
                p_nonce,
                p_qe_report,
                p_quote,
                maxlen,
                p_quote_len)
    };

    if result != sgx_status_t::SGX_SUCCESS {
        return Err(result);
    }

    if rt != sgx_status_t::SGX_SUCCESS {
        println!("ocall_get_quote returned {}", rt);
        return Err(rt);
    }

    // Perform a check on qe_report to verify if the qe_report is valid
    match tse::rsgx_verify_report(&qe_report) {
        Ok(()) => println!("rsgx_verify_report passed!"),
        Err(x) => {
            println!("rsgx_verify_report failed with {:?}", x);
            return Err(x);
        },
    }

    // Check if the qe_report is produced on the same platform
    if target_info.mr_enclave.m != qe_report.body.mr_enclave.m ||
       target_info.attributes.flags != qe_report.body.attributes.flags ||
       target_info.attributes.xfrm  != qe_report.body.attributes.xfrm {
        println!("qe_report does not match current target_info!");
        return Err(sgx_status_t::SGX_ERROR_UNEXPECTED);
    }

    println!("qe_report check passed");

    // Debug
    for i in 0..quote_len {
        print!("{:02X}", unsafe {*p_quote.offset(i as isize)});
    }
    println!("");

    // Check qe_report to defend against replay attack
    // The purpose of p_qe_report is for the ISV enclave to confirm the QUOTE
    // it received is not modified by the untrusted SW stack, and not a replay.
    // The implementation in QE is to generate a REPORT targeting the ISV
    // enclave (target info from p_report) , with the lower 32Bytes in
    // report.data = SHA256(p_nonce||p_quote). The ISV enclave can verify the
    // p_qe_report and report.data to confirm the QUOTE has not be modified and
    // is not a replay. It is optional.

    let mut rhs_vec : Vec<u8> = quote_nonce.rand.to_vec();
    rhs_vec.extend(&return_quote_buf[..quote_len as usize]);
    let rhs_hash = tcrypto::rsgx_sha256_slice(&rhs_vec[..]).unwrap();
    let lhs_hash = &qe_report.body.report_data.d[..32];

    //println!("rhs hash = {:02X}", rhs_hash.iter().format(""));
    //println!("report hs= {:02X}", lhs_hash.iter().format(""));

    if rhs_hash != lhs_hash {
        println!("Quote is tampered!");
        return Err(sgx_status_t::SGX_ERROR_UNEXPECTED);
    }

    let quote_vec : Vec<u8> = return_quote_buf[..quote_len as usize].to_vec();

    //let (attn_report, sig, cert) = get_report_from_intel(ias_sock, quote_vec);
    // Ok((attn_report, sig, cert))
    Err(sgx_status_t::SGX_ERROR_UNEXPECTED)
}
