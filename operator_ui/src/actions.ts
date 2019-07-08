import * as api from 'api'
import { AuthenticationError } from './api/errors'
import { pascalCase } from 'change-case'
import normalize from 'json-api-normalizer'

const createAction = type => ({ type: type })

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
    items: json.data.map(b => Object.assign({ id: b.id }, b.attributes))
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
    item: Object.assign({ id: json.data.id }, json.data.attributes)
  }),
  receiveErrorType: RECEIVE_BRIDGE_ERROR
}

function sendFetchActions(type, ...getArgs) {
  return dispatch => {
    const {
      requestActionType,
      receiveSuccess,
      receiveErrorType
    } = fetchActions[type]
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

const receiveSignInSuccess = json => {
  return {
    type: RECEIVE_SIGNIN_SUCCESS,
    authenticated: json.data.attributes.authenticated,
    errors: json.errors
  }
}

const receiveSignInFail = () => ({ type: RECEIVE_SIGNIN_FAIL })

function sendSignIn(data) {
  return dispatch => {
    dispatch(createAction(REQUEST_SIGNIN))
    return api
      .createSession(data)
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

function sendSignOut() {
  return dispatch => {
    dispatch(createAction(REQUEST_SIGNOUT))
    return api
      .destroySession()
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

export const fetchConfiguration = () => sendFetchActions('configuration')
export const fetchBridges = (page, size) =>
  sendFetchActions('bridges', page, size)
export const fetchBridgeSpec = name => sendFetchActions('bridgeSpec', name)

export const submitSignIn = data => sendSignIn(data)
export const submitSignOut = () => sendSignOut()

export const createJobSpec = (data, successCallback, errorCallback) => {
  return dispatch => {
    dispatch(createAction(REQUEST_CREATE))
    return api
      .createJobSpec(data)
      .then(res => {
        dispatch(receiveCreateSuccess(res))
        dispatch(notifySuccess(successCallback, res))
      })
      .catch(error => {
        curryErrorHandler(dispatch, RECEIVE_CREATE_ERROR)(error)
        dispatch(notifyError(errorCallback, error))
      })
  }
}

export const createJobRun = (id, successCallback, errorCallback) => {
  return dispatch => {
    dispatch(createAction(REQUEST_CREATE))
    return api
      .createJobSpecRun(id)
      .then(res => {
        dispatch(receiveCreateSuccess(res))
        dispatch(notifySuccess(successCallback, res))
      })
      .catch(error => {
        curryErrorHandler(dispatch, RECEIVE_CREATE_ERROR)(error)
        dispatch(notifyError(errorCallback, error))
      })
  }
}

export const createBridge = (data, successCallback, errorCallback) => {
  return dispatch => {
    dispatch(createAction(REQUEST_CREATE))
    return api
      .createBridge(data)
      .then(res => {
        dispatch(receiveCreateSuccess(res))
        dispatch(notifySuccess(successCallback, res.data))
      })
      .catch(error => {
        curryErrorHandler(dispatch, RECEIVE_CREATE_ERROR)(error)
        dispatch(notifyError(errorCallback, error))
      })
  }
}

export const updateBridge = (data, successCallback, errorCallback) => {
  return dispatch => {
    dispatch(createAction(REQUEST_UPDATE))
    return api
      .updateBridge(data)
      .then(res => {
        dispatch(receiveUpdateSuccess(res.data))
        dispatch(notifySuccess(successCallback, res.data))
      })
      .catch(error => {
        curryErrorHandler(dispatch, RECEIVE_UPDATE_ERROR)(error)
        dispatch(notifyError(errorCallback, error))
      })
  }
}

// DEV NOTE:
// Above here is deprecated. Use the `request(...)` function below to wrap API
// calls in a counter, normalize JSON-API responses and create notifications.
//
// The calls above will be converted gradually.
const handleError = dispatch => error => {
  if (error instanceof AuthenticationError) {
    dispatch(redirectToSignOut())
  } else {
    dispatch(notifyError(({ msg }) => msg, error))
  }
}

const request = (type, requestData, normalizeData, ...apiArgs) => {
  return dispatch => {
    dispatch({ type: `REQUEST_${type}` })
    return requestData(...apiArgs)
      .then(json => {
        const data = normalizeData(json)
        dispatch({ type: `UPSERT_${type}`, data: data })
      })
      .catch(handleError(dispatch))
      .finally(() => dispatch({ type: `RESPONSE_${type}` }))
  }
}

export const fetchAccountBalance = () =>
  request('ACCOUNT_BALANCE', api.getAccountBalance, json => normalize(json))

export const fetchJobs = (page, size) =>
  request(
    'JOBS',
    api.getJobs,
    json => normalize(json, { endpoint: 'currentPageJobs' }),
    page,
    size
  )

export const fetchRecentlyCreatedJobs = size =>
  request(
    'RECENTLY_CREATED_JOBS',
    api.getRecentlyCreatedJobs,
    json => normalize(json, { endpoint: 'recentlyCreatedJobs' }),
    size
  )

export const fetchJob = id =>
  request('JOB', api.getJobSpec, json => normalize(json), id)

export const fetchJobRuns = (opts: api.JobSpecRunsOpts) =>
  request(
    'JOB_RUNS',
    api.getJobSpecRuns,
    json => normalize(json, { endpoint: 'currentPageJobRuns' }),
    opts
  )

export const fetchRecentJobRuns = size =>
  request(
    'RECENT_JOB_RUNS',
    api.getRecentJobRuns,
    json => normalize(json, { endpoint: 'recentJobRuns' }),
    size
  )

export const fetchJobRun = id =>
  request('JOB_RUN', api.getJobSpecRun, json => normalize(json), id)

export const deleteCompletedJobRuns = updatedBefore =>
  request(
    'DELETE_COMPLETED_JOB_RUNS',
    api.bulkDeleteJobRuns,
    normalize,
    ['completed'],
    updatedBefore
  )

export const deleteErroredJobRuns = updatedBefore =>
  request(
    'DELETE_ERRORED_JOB_RUNS',
    api.bulkDeleteJobRuns,
    normalize,
    ['errored'],
    updatedBefore
  )

export const fetchTransactions = (page, size) =>
  request(
    'TRANSACTIONS',
    api.getTransactions,
    json => normalize(json, { endpoint: 'currentPageTransactions' }),
    page,
    size
  )

export const fetchTransaction = id =>
  request('TRANSACTION', api.getTransaction, json => normalize(json), id)
