import * as api from 'api'
import { AuthenticationError } from 'errors'
import { pascalCase } from 'change-case'

const fetchActions = {}

const createAction = (type) => ({type: type})

const requestNetworkError = (type, error) => ({
  type: type,
  error: error,
  networkError: true
})

export const MATCH_ROUTE = 'MATCH_ROUTE'

export const matchRoute = match => ({
  type: MATCH_ROUTE,
  match: match
})

export const REQUEST_JOBS = 'REQUEST_JOBS'
export const RECEIVE_JOBS_SUCCESS = 'RECEIVE_JOBS_SUCCESS'
export const RECEIVE_JOBS_ERROR = 'RECEIVE_JOBS_ERROR'

fetchActions.jobs = {
  action: createAction(REQUEST_JOBS),
  receiveSuccess: json => ({
    type: RECEIVE_JOBS_SUCCESS,
    count: json.meta.count,
    items: json.data.map((j) => (
      {
        id: j.id,
        createdAt: j.attributes.createdAt,
        initiators: j.attributes.initiators
      }
    ))
  }),
  receiveNetworkError: error => requestNetworkError(RECEIVE_JOBS_ERROR, error)
}

export const REQUEST_ACCOUNT_BALANCE = 'REQUEST_ACCOUNT_BALANCE'
export const RECEIVE_ACCOUNT_BALANCE_SUCCESS = 'RECEIVE_ACCOUNT_BALANCE_SUCCESS'
export const RECEIVE_ACCOUNT_BALANCE_ERROR = 'RECEIVE_ACCOUNT_BALANCE_ERROR'

fetchActions.accountBalance = {
  action: createAction(REQUEST_ACCOUNT_BALANCE),
  receiveSuccess: json => ({
    type: RECEIVE_ACCOUNT_BALANCE_SUCCESS,
    eth: json.data.attributes.ethBalance,
    link: json.data.attributes.linkBalance
  }),
  receiveNetworkError: error => requestNetworkError(RECEIVE_ACCOUNT_BALANCE_ERROR, error)
}

export const REQUEST_JOB_SPEC = 'REQUEST_JOB_SPEC'
export const RECEIVE_JOB_SPEC_SUCCESS = 'RECEIVE_JOB_SPEC_SUCCESS'
export const RECEIVE_JOB_SPEC_ERROR = 'RECEIVE_JOB_SPEC_ERROR'

fetchActions.jobSpec = {
  action: createAction(REQUEST_JOB_SPEC),
  receiveSuccess: json => ({
    type: RECEIVE_JOB_SPEC_SUCCESS,
    item: json.data.attributes
  }),
  receiveNetworkError: error => requestNetworkError(RECEIVE_JOB_SPEC_ERROR, error)
}

export const REQUEST_JOB_SPEC_RUNS = 'REQUEST_JOB_SPEC_RUNS'
export const RECEIVE_JOB_SPEC_RUNS_SUCCESS = 'RECEIVE_JOB_SPEC_RUNS_SUCCESS'
export const RECEIVE_JOB_SPEC_RUNS_ERROR = 'RECEIVE_JOB_SPEC_RUNS_ERROR'

fetchActions.jobSpecRuns = {
  action: createAction(REQUEST_JOB_SPEC_RUNS),
  receiveSuccess: json => ({
    type: RECEIVE_JOB_SPEC_RUNS_SUCCESS,
    items: json.data.map(j => j.attributes),
    runsCount: json.meta.count
  }),
  receiveNetworkError: error => requestNetworkError(RECEIVE_JOB_SPEC_RUNS_ERROR, error)
}

export const REQUEST_JOB_SPEC_RUN = 'REQUEST_JOB_SPEC_RUN'
export const RECEIVE_JOB_SPEC_RUN_SUCCESS = 'RECEIVE_JOB_SPEC_RUN_SUCCESS'
export const RECEIVE_JOB_SPEC_RUN_ERROR = 'RECEIVE_JOB_SPEC_RUN_ERROR'

fetchActions.jobSpecRun = {
  action: createAction(REQUEST_JOB_SPEC_RUN),
  receiveSuccess: json => ({
    type: RECEIVE_JOB_SPEC_RUN_SUCCESS,
    item: json.data.attributes
  }),
  receiveNetworkError: error => requestNetworkError(RECEIVE_JOB_SPEC_RUN_ERROR, error)
}

export const REQUEST_CONFIGURATION = 'REQUEST_CONFIGURATION'
export const RECEIVE_CONFIGURATION_SUCCESS = 'RECEIVE_CONFIGURATION_SUCCESS'
export const RECEIVE_CONFIGURATION_ERROR = 'RECEIVE_CONFIGURATION_ERROR'

fetchActions.configuration = {
  action: createAction(REQUEST_CONFIGURATION),
  receiveSuccess: json => ({
    type: RECEIVE_CONFIGURATION_SUCCESS,
    config: json.data.attributes
  }),
  receiveNetworkError: error => requestNetworkError(RECEIVE_CONFIGURATION_ERROR, error)
}

export const REQUEST_BRIDGES = 'REQUEST_BRIDGES'
export const RECEIVE_BRIDGES_SUCCESS = 'RECEIVE_BRIDGES_SUCCESS'
export const RECEIVE_BRIDGES_ERROR = 'RECEIVE_BRIDGES_ERROR'

fetchActions.bridges = {
  action: createAction(REQUEST_BRIDGES),
  receiveSuccess: json => ({
    type: RECEIVE_BRIDGES_SUCCESS,
    count: json.meta.count,
    items: json.data.map(b => b.attributes)
  }),
  receiveNetworkError: error => requestNetworkError(RECEIVE_BRIDGES_ERROR, error)
}

export const REQUEST_BRIDGESPEC = 'REQUEST_BRIDGESPEC'
export const RECEIVE_BRIDGESPEC_SUCCESS = 'RECEIVE_BRIDGESPEC_SUCCESS'
export const RECEIVE_BRIDGESPEC_ERROR = 'RECEIVE_BRIDGESPEC_ERROR'

fetchActions.bridgeSpec = {
  action: createAction(REQUEST_BRIDGESPEC),
  receiveSuccess: json => ({
    type: RECEIVE_BRIDGESPEC_SUCCESS,
    name: json.data.attributes.name,
    url: json.data.attributes.url,
    confirmations: json.data.attributes.defaultConfirmations,
    incomingToken: json.data.attributes.incomingToken,
    outgoingToken: json.data.attributes.outgoingToken
  }),
  receiveNetworkError: error => requestNetworkError(RECEIVE_BRIDGESPEC_ERROR, error)
}

function sendFetchActions (type, ...getArgs) {
  return dispatch => {
    const {action, receiveSuccess, receiveNetworkError} = fetchActions[type]
    const apiGet = api['get' + pascalCase(type)]

    dispatch(action)
    return apiGet(...getArgs)
      .then(json => dispatch(receiveSuccess(json)))
      .catch(error => dispatch(receiveNetworkError(error)))
  }
}

export const REQUEST_SIGNIN = 'REQUEST_SIGNIN'
export const RECEIVE_SIGNIN_SUCCESS = 'RECEIVE_SIGNIN_SUCCESS'
export const RECEIVE_SIGNIN_FAIL = 'RECEIVE_SIGNIN_FAIL'
export const RECEIVE_SIGNIN_ERROR = 'RECEIVE_SIGNIN_ERROR'

const receiveSignInSuccess = (json) => ({
  type: RECEIVE_SIGNIN_SUCCESS,
  authenticated: json.authenticated,
  errors: json.errors
})

const receiveSignInFail = () => ({type: RECEIVE_SIGNIN_FAIL})

function sendSignIn (data) {
  return dispatch => {
    dispatch(createAction(REQUEST_SIGNIN))
    return api.createSession(data)
      .then((json) => dispatch(receiveSignInSuccess(json)))
      .catch(error => {
        if (error instanceof AuthenticationError) {
          dispatch(receiveSignInFail())
        } else {
          dispatch(requestNetworkError(RECEIVE_SIGNIN_ERROR, error))
        }
      })
  }
}

export const REQUEST_SIGNOUT = 'REQUEST_SIGNOUT'
export const RECEIVE_SIGNOUT_SUCCESS = 'RECEIVE_SIGNOUT_SUCCESS'
export const RECEIVE_SIGNOUT_ERROR = 'RECEIVE_SIGNOUT_ERROR'

export const receiveSignoutSuccess = () => ({
  type: RECEIVE_SIGNOUT_SUCCESS,
  authenticated: false
})

function sendSignOut () {
  return dispatch => {
    dispatch(createAction(REQUEST_SIGNOUT))
    return api.destroySession()
      .then((json) => dispatch(receiveSignoutSuccess(json)))
      .catch(error => dispatch(requestNetworkError(RECEIVE_SIGNOUT_ERROR, error)))
  }
}

export const REQUEST_CREATE = 'REQUEST_CREATE'
export const RECEIVE_CREATE_SUCCESS = 'RECEIVE_CREATE_SUCCESS'
export const RECEIVE_CREATE_ERROR = 'RECEIVE_CREATE_ERROR'

const receiveCreateSuccess = (res) => ({
  type: RECEIVE_CREATE_SUCCESS,
  response: res
})

function sendJobSpec (data, shouldStringify) {
  return dispatch => {
    dispatch(createAction(REQUEST_CREATE))
    return api.createJobSpec(data, shouldStringify)
      .then((res) => dispatch(receiveCreateSuccess(res)))
      .catch(error => dispatch(requestNetworkError(RECEIVE_CREATE_ERROR, error)))
  }
}

function sendBridgeType (data, shouldStringify) {
  return dispatch => {
    dispatch(createAction(REQUEST_CREATE))
    return api.createBridgeType(data, shouldStringify)
      .then((res) => dispatch(receiveCreateSuccess(res)))
      .catch(error => dispatch(requestNetworkError(RECEIVE_CREATE_ERROR, error)))
  }
}

function sendJobSpecRun (id) {
  return dispatch => {
    dispatch(createAction(REQUEST_CREATE))
    return api.createJobSpecRun(id)
      .then((res) => dispatch(receiveCreateSuccess(res)))
      .catch((error) => dispatch(requestNetworkError(RECEIVE_CREATE_ERROR, error)))
  }
}

export const fetchJobs = (page, size) => sendFetchActions('jobs', page, size)
export const fetchAccountBalance = () => sendFetchActions('accountBalance')
export const fetchJobSpec = (id) => sendFetchActions('jobSpec', id)
export const fetchJobSpecRuns = (id, page, size) => sendFetchActions('jobSpecRuns', id, page, size)
export const fetchJobSpecRun = (id) => sendFetchActions('jobSpecRun', id)
export const fetchConfiguration = () => sendFetchActions('configuration')
export const fetchBridges = (page, size) => sendFetchActions('bridges', page, size)
export const fetchBridgeSpec = (name) => sendFetchActions('bridgeSpec', name)

export const submitSignIn = (data) => sendSignIn(data)
export const submitSignOut = () => sendSignOut()
export const submitBridgeType = (data, shouldStringify) => sendBridgeType(data, shouldStringify)
export const submitJobSpec = (data, shouldStringify) => sendJobSpec(data, shouldStringify)
export const submitJobSpecRun = (id) => sendJobSpecRun(id)
