import {
  getAccountBalance,
  getConfiguration,
  getJobs,
  getJobSpec,
  getJobSpecRuns
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
    item: json.data.attributes
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

export const RECEIVE_JOB_SPEC_RUNS_SUCCESS = 'RECEIVE_JOB_SPEC_RUNS_SUCCESS'
export const RECEIVE_JOB_SPEC_RUNS_ERROR = 'RECEIVE_JOB_SPEC_RUNS_ERROR'

const receiveJobSpecRunsSuccess = (json) => {
  return {
    type: RECEIVE_JOB_SPEC_RUNS_SUCCESS,
    items: json.data.map(j => j.attributes),
    runsCount: json.meta.count
  }
}
const receiveJobSpecRunsNetworkError = () => {
  return {
    type: RECEIVE_JOB_SPEC_RUNS_ERROR,
    networkError: true
  }
}

export const fetchJobSpecRuns = (id, page, size) => {
  return dispatch => {
    return getJobSpecRuns(id, page, size)
      .then(json => dispatch(receiveJobSpecRunsSuccess(json)))
      .catch(_ => dispatch(receiveJobSpecRunsNetworkError()))
  }
}

export const REQUEST_CONFIGURATION = 'REQUEST_CONFIGURATION'
export const RECEIVE_CONFIGURATION_SUCCESS = 'RECEIVE_CONFIGURATION_SUCCESS'
export const RECEIVE_CONFIGURATION_ERROR = 'RECEIVE_CONFIGURATION_ERROR'

const requestConfiguration = () => ({ type: REQUEST_CONFIGURATION })
const receiveConfiguration = (json) => {
  return {
    type: RECEIVE_CONFIGURATION_SUCCESS,
    config: json.data.attributes
  }
}
const receiveConfigurationNetworkError = () => {
  return {
    type: RECEIVE_CONFIGURATION_ERROR,
    networkError: true
  }
}

export const fetchConfiguration = () => {
  return dispatch => {
    dispatch(requestConfiguration())
    return getConfiguration()
      .then(json => dispatch(receiveConfiguration(json)))
      .catch(_ => dispatch(receiveConfigurationNetworkError()))
  }
}
