import {
  getAccountBalance,
  getJobs,
  getJobSpec
} from 'api'

export const REQUEST_JOBS = 'REQUEST_JOBS'
export const RECEIVE_JOBS_SUCCESS = 'RECEIVE_JOBS_SUCCESS'
export const RECEIVE_JOBS_ERROR = 'RECEIVE_JOBS_ERROR'

const requestJobs = () => ({ type: REQUEST_JOBS })
const receiveJobsSuccess = (json) => {
  return {
    type: RECEIVE_JOBS_SUCCESS,
    count: json.meta.count,
    items: json.data.map((j) => (
      {
        id: j.id,
        createdAt: j.attributes.createdAt,
        initiators: j.attributes.initiators
      }
    ))
  }
}
const receiveJobsNetworkError = () => {
  return {
    type: RECEIVE_JOBS_ERROR,
    networkError: true
  }
}

export const fetchJobs = (page, size) => {
  return dispatch => {
    dispatch(requestJobs())
    return getJobs(page, size)
      .then(json => dispatch(receiveJobsSuccess(json)))
      .catch(_ => dispatch(receiveJobsNetworkError()))
  }
}

export const REQUEST_ACCOUNT_BALANCE = 'REQUEST_ACCOUNT_BALANCE'
export const RECEIVE_ACCOUNT_BALANCE_SUCCESS = 'RECEIVE_ACCOUNT_BALANCE_SUCCESS'
export const RECEIVE_ACCOUNT_BALANCE_ERROR = 'RECEIVE_ACCOUNT_BALANCE_ERROR'

const requestAccountBalance = () => ({ type: REQUEST_ACCOUNT_BALANCE })
const receiveAccountBalance = (json) => {
  return {
    type: RECEIVE_ACCOUNT_BALANCE_SUCCESS,
    eth: json.data.attributes.ethBalance,
    link: json.data.attributes.linkBalance
  }
}
const receiveAccountBalanceNetworkError = () => {
  return {
    type: RECEIVE_ACCOUNT_BALANCE_ERROR,
    networkError: true
  }
}

export const fetchAccountBalance = () => {
  return dispatch => {
    dispatch(requestAccountBalance())
    return getAccountBalance()
      .then(json => dispatch(receiveAccountBalance(json)))
      .catch(_ => dispatch(receiveAccountBalanceNetworkError()))
  }
}

export const REQUEST_JOB_SPEC = 'REQUEST_JOB_SPEC'
export const RECEIVE_JOB_SPEC_SUCCESS = 'RECEIVE_JOB_SPEC_SUCCESS'
export const RECEIVE_JOB_SPEC_ERROR = 'RECEIVE_JOB_SPEC_ERROR'

const requestJobSpec = () => ({ type: REQUEST_JOB_SPEC })
const receiveJobSpecSuccess = (json) => {
  return {
    type: RECEIVE_JOB_SPEC_SUCCESS,
    item: json
  }
}
const receiveJobSpecNetworkError = () => {
  return {
    type: RECEIVE_JOB_SPEC_ERROR,
    networkError: true
  }
}

export const fetchJobSpec = (id) => {
  return dispatch => {
    dispatch(requestJobSpec())
    return getJobSpec(id)
      .then(json => dispatch(receiveJobSpecSuccess(json)))
      .catch(_ => dispatch(receiveJobSpecNetworkError()))
  }
}
