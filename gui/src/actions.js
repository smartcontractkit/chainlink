import * as api from 'api'
import { AuthenticationError } from 'errors'
import { pascalCase } from 'change-case'
import transformJobs from 'actions/transforms/jobs'

const createAction = type => ({type: type})

const createErrorAction = (error, type) => ({
  type: type,
  error: error.stack,
  networkError: true
})

const curryErrorHandler = (dispatch, type) => error => {
  if (error instanceof AuthenticationError) {
    dispatch(redirectToSignOut())
  } else {
    dispatch(createErrorAction(error, type))
  }
}

export const REDIRECT = 'REDIRECT'

const redirectToSignOut = () => ({
  type: REDIRECT,
  to: '/signout'
})

export const MATCH_ROUTE = 'MATCH_ROUTE'

export const matchRoute = match => ({
  type: MATCH_ROUTE,
  match: match
})

export const NOTIFY_SUCCESS = 'NOTIFY_SUCCESS'

export const notifySuccess = (component, props) => ({
  type: NOTIFY_SUCCESS,
  component: component,
  props: props
})

export const NOTIFY_ERROR = 'NOTIFY_ERROR'

export const notifyError = (component, error) => ({
  type: NOTIFY_ERROR,
  component: component,
  error: error
})

const fetchActions = {}

export const REQUEST_JOBS = 'REQUEST_JOBS'
export const RECEIVE_JOBS_SUCCESS = 'RECEIVE_JOBS_SUCCESS'
export const RECEIVE_JOBS_ERROR = 'RECEIVE_JOBS_ERROR'

fetchActions.jobs = {
  requestActionType: REQUEST_JOBS,
  receiveSuccess: json => transformJobs(RECEIVE_JOBS_SUCCESS, json),
  receiveErrorType: RECEIVE_JOBS_ERROR
}

export const REQUEST_RECENTLY_CREATED_JOBS = 'REQUEST_RECENTLY_CREATED_JOBS'
export const RECEIVE_RECENTLY_CREATED_JOBS_SUCCESS = 'RECEIVE_RECENTLY_CREATED_JOBS_SUCCESS'
export const RECEIVE_RECENTLY_CREATED_JOBS_ERROR = 'RECEIVE_RECENTLY_CREATED_JOBS_ERROR'

fetchActions.recentlyCreatedJobs = {
  requestActionType: REQUEST_RECENTLY_CREATED_JOBS,
  receiveSuccess: json => transformJobs(RECEIVE_RECENTLY_CREATED_JOBS_SUCCESS, json),
  receiveErrorType: RECEIVE_RECENTLY_CREATED_JOBS_ERROR
}

export const REQUEST_ACCOUNT_BALANCE = 'REQUEST_ACCOUNT_BALANCE'
export const RECEIVE_ACCOUNT_BALANCE_SUCCESS = 'RECEIVE_ACCOUNT_BALANCE_SUCCESS'
export const RECEIVE_ACCOUNT_BALANCE_ERROR = 'RECEIVE_ACCOUNT_BALANCE_ERROR'

fetchActions.accountBalance = {
  requestActionType: REQUEST_ACCOUNT_BALANCE,
  receiveSuccess: json => ({
    type: RECEIVE_ACCOUNT_BALANCE_SUCCESS,
    eth: json.data.attributes.ethBalance,
    link: json.data.attributes.linkBalance
  }),
  receiveErrorType: RECEIVE_ACCOUNT_BALANCE_ERROR
}

export const REQUEST_JOB_SPEC = 'REQUEST_JOB_SPEC'
export const RECEIVE_JOB_SPEC_SUCCESS = 'RECEIVE_JOB_SPEC_SUCCESS'
export const RECEIVE_JOB_SPEC_ERROR = 'RECEIVE_JOB_SPEC_ERROR'

fetchActions.jobSpec = {
  requestActionType: REQUEST_JOB_SPEC,
  receiveSuccess: json => ({
    type: RECEIVE_JOB_SPEC_SUCCESS,
    item: json.data.attributes
  }),
  receiveErrorType: RECEIVE_JOB_SPEC_ERROR
}

export const REQUEST_JOB_SPEC_RUNS = 'REQUEST_JOB_SPEC_RUNS'
export const RECEIVE_JOB_SPEC_RUNS_SUCCESS = 'RECEIVE_JOB_SPEC_RUNS_SUCCESS'
export const RECEIVE_JOB_SPEC_RUNS_ERROR = 'RECEIVE_JOB_SPEC_RUNS_ERROR'

fetchActions.jobSpecRuns = {
  requestActionType: REQUEST_JOB_SPEC_RUNS,
  receiveSuccess: json => ({
    type: RECEIVE_JOB_SPEC_RUNS_SUCCESS,
    items: json.data.map(j => j.attributes),
    runsCount: json.meta.count
  }),
  receiveErrorType: RECEIVE_JOB_SPEC_RUNS_ERROR
}

export const REQUEST_JOB_SPEC_RUN = 'REQUEST_JOB_SPEC_RUN'
export const RECEIVE_JOB_SPEC_RUN_SUCCESS = 'RECEIVE_JOB_SPEC_RUN_SUCCESS'
export const RECEIVE_JOB_SPEC_RUN_ERROR = 'RECEIVE_JOB_SPEC_RUN_ERROR'

fetchActions.jobSpecRun = {
  requestActionType: REQUEST_JOB_SPEC_RUN,
  receiveSuccess: json => ({
    type: RECEIVE_JOB_SPEC_RUN_SUCCESS,
    item: json.data.attributes
  }),
  receiveErrorType: RECEIVE_JOB_SPEC_RUN_ERROR
}

export const REQUEST_CONFIGURATION = 'REQUEST_CONFIGURATION'
export const RECEIVE_CONFIGURATION_SUCCESS = 'RECEIVE_CONFIGURATION_SUCCESS'
export const RECEIVE_CONFIGURATION_ERROR = 'RECEIVE_CONFIGURATION_ERROR'

fetchActions.configuration = {
  requestActionType: REQUEST_CONFIGURATION,
  receiveSuccess: json => ({
    type: RECEIVE_CONFIGURATION_SUCCESS,
    config: json.data.attributes
  }),
  receiveErrorType: RECEIVE_CONFIGURATION_ERROR
}

export const REQUEST_BRIDGES = 'REQUEST_BRIDGES'
export const RECEIVE_BRIDGES_SUCCESS = 'RECEIVE_BRIDGES_SUCCESS'
export const RECEIVE_BRIDGES_ERROR = 'RECEIVE_BRIDGES_ERROR'

fetchActions.bridges = {
  requestActionType: REQUEST_BRIDGES,
  receiveSuccess: json => ({
    type: RECEIVE_BRIDGES_SUCCESS,
    count: json.meta.count,
    items: json.data.map(b => Object.assign({id: b.id}, b.attributes))
  }),
  receiveErrorType: RECEIVE_BRIDGES_ERROR
}

export const REQUEST_BRIDGE = 'REQUEST_BRIDGE'
export const RECEIVE_BRIDGE_SUCCESS = 'RECEIVE_BRIDGE_SUCCESS'
export const RECEIVE_BRIDGE_ERROR = 'RECEIVE_BRIDGE_ERROR'

fetchActions.bridgeSpec = {
  requestActionType: REQUEST_BRIDGE,
  receiveSuccess: json => ({
    type: RECEIVE_BRIDGE_SUCCESS,
    item: Object.assign(
      {id: json.data.id},
      json.data.attributes
    )
  }),
  receiveErrorType: RECEIVE_BRIDGE_ERROR
}

function sendFetchActions (type, ...getArgs) {
  return dispatch => {
    const {requestActionType, receiveSuccess, receiveErrorType} = fetchActions[type]
    const apiGet = api['get' + pascalCase(type)]

    dispatch(createAction(requestActionType))
    return apiGet(...getArgs)
      .then(json => dispatch(receiveSuccess(json)))
      .catch(curryErrorHandler(dispatch, receiveErrorType))
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
      .then(json => dispatch(receiveSignInSuccess(json)))
      .catch(error => {
        if (error instanceof AuthenticationError) {
          dispatch(receiveSignInFail())
        } else {
          dispatch(createErrorAction(error, RECEIVE_SIGNIN_ERROR))
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
      .then(json => dispatch(receiveSignoutSuccess(json)))
      .catch(curryErrorHandler(dispatch, RECEIVE_SIGNIN_ERROR))
  }
}

export const REQUEST_CREATE = 'REQUEST_CREATE'
export const RECEIVE_CREATE_SUCCESS = 'RECEIVE_CREATE_SUCCESS'
export const RECEIVE_CREATE_ERROR = 'RECEIVE_CREATE_ERROR'

const receiveCreateSuccess = response => ({
  type: RECEIVE_CREATE_SUCCESS,
  response: response
})

export const REQUEST_UPDATE = 'REQUEST_UPDATE'
export const RECEIVE_UPDATE_SUCCESS = 'RECEIVE_UPDATE_SUCCESS'
export const RECEIVE_UPDATE_ERROR = 'RECEIVE_UPDATE_ERROR'

const receiveUpdateSuccess = response => ({
  type: RECEIVE_UPDATE_SUCCESS,
  response: response
})

export const fetchJobs = (page, size) => sendFetchActions('jobs', page, size)
export const fetchRecentlyCreatedJobs = size => sendFetchActions('recentlyCreatedJobs', size)
export const fetchAccountBalance = () => sendFetchActions('accountBalance')
export const fetchJobSpec = id => sendFetchActions('jobSpec', id)
export const fetchJobSpecRuns = (id, page, size) => sendFetchActions('jobSpecRuns', id, page, size)
export const fetchJobSpecRun = id => sendFetchActions('jobSpecRun', id)
export const fetchConfiguration = () => sendFetchActions('configuration')
export const fetchBridges = (page, size) => sendFetchActions('bridges', page, size)
export const fetchBridgeSpec = name => sendFetchActions('bridgeSpec', name)

export const submitSignIn = data => sendSignIn(data)
export const submitSignOut = () => sendSignOut()

export const createJobSpec = (data, renderNotifySuccess, renderNotifyError) => {
  return dispatch => {
    dispatch(createAction(REQUEST_CREATE))
    return api.createJobSpec(data)
      .then(res => {
        dispatch(receiveCreateSuccess(res))
        dispatch(notifySuccess(renderNotifySuccess, res))
      })
      .catch(error => {
        curryErrorHandler(dispatch, RECEIVE_CREATE_ERROR)(error)
        dispatch(notifyError(renderNotifyError, error))
      })
  }
}

export const createJobRun = (id, renderNotifySuccess, renderNotifyError) => {
  return dispatch => {
    dispatch(createAction(REQUEST_CREATE))
    return api.createJobSpecRun(id)
      .then(res => {
        dispatch(receiveCreateSuccess(res))
        dispatch(notifySuccess(renderNotifySuccess, res))
      })
      .catch(error => {
        curryErrorHandler(dispatch, RECEIVE_CREATE_ERROR)(error)
        dispatch(notifyError(renderNotifyError, error))
      })
  }
}

export const createBridge = (data, renderNotifySuccess, renderNotifyError) => {
  return dispatch => {
    dispatch(createAction(REQUEST_CREATE))
    return api.createBridge(data)
      .then(res => {
        dispatch(receiveCreateSuccess(res))
        dispatch(notifySuccess(renderNotifySuccess, data))
      })
      .catch(error => {
        curryErrorHandler(dispatch, RECEIVE_CREATE_ERROR)(error)
        dispatch(notifyError(renderNotifyError, error))
      })
  }
}

export const updateBridge = (data, renderNotifySuccess, renderNotifyError) => {
  return dispatch => {
    dispatch(createAction(REQUEST_UPDATE))
    return api.updateBridge(data)
      .then(res => {
        dispatch(receiveUpdateSuccess(res))
        dispatch(notifySuccess(renderNotifySuccess, data))
      })
      .catch(error => {
        curryErrorHandler(dispatch, RECEIVE_UPDATE_ERROR)(error)
        dispatch(notifyError(renderNotifyError, error))
      })
  }
}
