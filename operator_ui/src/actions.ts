import * as api from './api'
import { Dispatch } from 'redux'
import { AuthenticationError } from './api/errors'
import { pascalCase } from 'change-case'
import normalize from 'json-api-normalizer'

const createAction = (type: string) => ({ type: type })

const createErrorAction = (error: Error, type: string) => ({
  type: type,
  error: error.stack,
  networkError: true
})

const curryErrorHandler = (dispatch: Dispatch, type: string) => (
  error: Error
) => {
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

export const matchRoute = (match: any) => ({
  type: MATCH_ROUTE,
  match: match
})

export const NOTIFY_SUCCESS = 'NOTIFY_SUCCESS'

export const notifySuccess = (component: any, props: any) => ({
  type: NOTIFY_SUCCESS,
  component: component,
  props: props
})

export const NOTIFY_ERROR = 'NOTIFY_ERROR'

export const notifyError = (component: any, error: Error) => ({
  type: NOTIFY_ERROR,
  component: component,
  error: error
})

// TODO: These should no longer be needed after refactor
export const REQUEST_CONFIGURATION = 'REQUEST_CONFIGURATION'
export const RECEIVE_CONFIGURATION_SUCCESS = 'RECEIVE_CONFIGURATION_SUCCESS'
export const RECEIVE_CONFIGURATION_ERROR = 'RECEIVE_CONFIGURATION_ERROR'

export const REQUEST_BRIDGES = 'REQUEST_BRIDGES'
export const RECEIVE_BRIDGES_SUCCESS = 'RECEIVE_BRIDGES_SUCCESS'
export const RECEIVE_BRIDGES_ERROR = 'RECEIVE_BRIDGES_ERROR'

export const REQUEST_BRIDGE = 'REQUEST_BRIDGE'
export const RECEIVE_BRIDGE_SUCCESS = 'RECEIVE_BRIDGE_SUCCESS'
export const RECEIVE_BRIDGE_ERROR = 'RECEIVE_BRIDGE_ERROR'
// END TODO: These should no longer be needed after refactor

export const REQUEST_SIGNIN = 'REQUEST_SIGNIN'
export const RECEIVE_SIGNIN_SUCCESS = 'RECEIVE_SIGNIN_SUCCESS'
export const RECEIVE_SIGNIN_FAIL = 'RECEIVE_SIGNIN_FAIL'
export const RECEIVE_SIGNIN_ERROR = 'RECEIVE_SIGNIN_ERROR'

const receiveSignInSuccess = (json: any) => {
  return {
    type: RECEIVE_SIGNIN_SUCCESS,
    authenticated: json.data.attributes.authenticated,
    errors: json.errors
  }
}

const receiveSignInFail = () => ({ type: RECEIVE_SIGNIN_FAIL })

function sendSignIn(data: object) {
  return (dispatch: Dispatch) => {
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
  return (dispatch: Dispatch) => {
    dispatch(createAction(REQUEST_SIGNOUT))
    return api
      .destroySession()
      .then(() => dispatch(receiveSignoutSuccess()))
      .catch(curryErrorHandler(dispatch, RECEIVE_SIGNIN_ERROR))
  }
}

export const REQUEST_CREATE = 'REQUEST_CREATE'
export const RECEIVE_CREATE_SUCCESS = 'RECEIVE_CREATE_SUCCESS'
export const RECEIVE_CREATE_ERROR = 'RECEIVE_CREATE_ERROR'

const receiveCreateSuccess = (response: Response) => ({
  type: RECEIVE_CREATE_SUCCESS,
  response: response
})

export const REQUEST_DELETE = 'REQUEST_DELETE'
export const RECEIVE_DELETE_SUCCESS = 'RECEIVE_DELETE_SUCCESS'
export const RECEIVE_DELETE_ERROR = 'RECEIVE_DELETE_ERROR'

const receiveDeleteSuccess = (response: Response) => ({
  type: RECEIVE_DELETE_SUCCESS,
  response: response
})

export const REQUEST_UPDATE = 'REQUEST_UPDATE'
export const RECEIVE_UPDATE_SUCCESS = 'RECEIVE_UPDATE_SUCCESS'
export const RECEIVE_UPDATE_ERROR = 'RECEIVE_UPDATE_ERROR'

const receiveUpdateSuccess = (response: Response) => ({
  type: RECEIVE_UPDATE_SUCCESS,
  response: response
})

export const submitSignIn = (data: object) => sendSignIn(data)
export const submitSignOut = () => sendSignOut()

export const createJobSpec = (
  data: any,
  successCallback: any,
  errorCallback: any
) => {
  return (dispatch: Dispatch) => {
    dispatch(createAction(REQUEST_CREATE))
    return api
      .createJobSpec(data)
      .then((res: any) => {
        dispatch(receiveCreateSuccess(res))
        dispatch(notifySuccess(successCallback, res))
      })
      .catch(error => {
        curryErrorHandler(dispatch, RECEIVE_CREATE_ERROR)(error)
        dispatch(notifyError(errorCallback, error))
      })
  }
}

export const deleteJobSpec = (
  data: any,
  successCallback: any,
  errorCallback: any
) => {
  return (dispatch: Dispatch) => {
    dispatch(createAction(REQUEST_DELETE))
    return api
      .destroyJobSpec(data)
      .then(res => {
        dispatch(receiveDeleteSuccess(data))
        dispatch(notifySuccess(successCallback, res))
      })
      .catch(error => {
        curryErrorHandler(dispatch, RECEIVE_DELETE_ERROR)(error)
        dispatch(notifyError(errorCallback, error))
      })
  }
}

export const createJobRun = (
  id: string,
  successCallback: any,
  errorCallback: any
) => {
  return (dispatch: Dispatch) => {
    dispatch(createAction(REQUEST_CREATE))
    return api
      .createJobSpecRun(id)
      .then((res: any) => {
        dispatch(receiveCreateSuccess(res))
        dispatch(notifySuccess(successCallback, res))
      })
      .catch(error => {
        curryErrorHandler(dispatch, RECEIVE_CREATE_ERROR)(error)
        dispatch(notifyError(errorCallback, error))
      })
  }
}

export const createBridge = (
  data: object,
  successCallback: any,
  errorCallback: any
) => {
  return (dispatch: Dispatch) => {
    dispatch(createAction(REQUEST_CREATE))
    return api
      .createBridge(data)
      .then((res: any) => {
        dispatch(receiveCreateSuccess(res))
        dispatch(notifySuccess(successCallback, res.data))
      })
      .catch(error => {
        curryErrorHandler(dispatch, RECEIVE_CREATE_ERROR)(error)
        dispatch(notifyError(errorCallback, error))
      })
  }
}

export const updateBridge = (
  data: object,
  successCallback: any,
  errorCallback: any
) => {
  return (dispatch: Dispatch) => {
    dispatch(createAction(REQUEST_UPDATE))
    return api
      .updateBridge(data)
      .then((res: any) => {
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
const handleError = (dispatch: Dispatch) => (error: Error) => {
  if (error instanceof AuthenticationError) {
    dispatch(redirectToSignOut())
  } else {
    dispatch(notifyError(({ msg }: any) => msg, error))
  }
}

const request = (
  type: string,
  requestData: any,
  normalizeData: any,
  ...apiArgs: any
) => {
  return (dispatch: Dispatch) => {
    dispatch({ type: `REQUEST_${type}` })
    return requestData(...apiArgs)
      .then((json: object) => {
        const data = normalizeData(json)
        dispatch({ type: `UPSERT_${type}`, data: data })
      })
      .catch(handleError(dispatch))
      .finally(() => dispatch({ type: `RESPONSE_${type}` }))
  }
}

export const fetchAccountBalance = () =>
  request('ACCOUNT_BALANCE', api.getAccountBalance, (json: any) =>
    normalize(json)
  )

// TODO: Reducers need to now work with UPSERT_CONFIGURATION & RESPONSE_CONFIGURATION
export const fetchConfiguration = () =>
  request('CONFIGURATION', api.getConfiguration, (json: any) => normalize(json))

// TODO: Reducers need to now work with UPSERT_BRIDGES & RESPONSE_BRIDGES
export const fetchBridges = (page: number, size: number) =>
  request('BRIDGES', api.getBridges, (json: any) => normalize(json), page, size)

// TODO: Reducers need to now work with UPSERT_BRIDGE & RESPONSE_BRIDGE
export const fetchBridgeSpec = (name: string) =>
  request('BRIDGE', api.getBridgeSpec, (json: any) => normalize(json), name)

export const fetchJobs = (page: number, size: number) =>
  request(
    'JOBS',
    api.getJobs,
    (json: any) => normalize(json, { endpoint: 'currentPageJobs' }),
    page,
    size
  )

export const fetchRecentlyCreatedJobs = (size: number) =>
  request(
    'RECENTLY_CREATED_JOBS',
    api.getRecentlyCreatedJobs,
    (json: any) => normalize(json, { endpoint: 'recentlyCreatedJobs' }),
    size
  )

export const fetchJob = (id: string) =>
  request('JOB', api.getJobSpec, (json: object) => normalize(json), id)

export const fetchJobRuns = (opts: api.JobSpecRunsOpts) =>
  request(
    'JOB_RUNS',
    api.getJobSpecRuns,
    (json: any) => normalize(json, { endpoint: 'currentPageJobRuns' }),
    opts
  )

export const fetchRecentJobRuns = (size: number) =>
  request(
    'RECENT_JOB_RUNS',
    api.getRecentJobRuns,
    (json: any) => normalize(json, { endpoint: 'recentJobRuns' }),
    size
  )

export const fetchJobRun = (id: string) =>
  request('JOB_RUN', api.getJobSpecRun, (json: any) => normalize(json), id)

export const deleteCompletedJobRuns = (updatedBefore: any) =>
  request(
    'DELETE_COMPLETED_JOB_RUNS',
    api.bulkDeleteJobRuns,
    normalize,
    ['completed'],
    updatedBefore
  )

export const deleteErroredJobRuns = (updatedBefore: any) =>
  request(
    'DELETE_ERRORED_JOB_RUNS',
    api.bulkDeleteJobRuns,
    normalize,
    ['errored'],
    updatedBefore
  )

export const fetchTransactions = (page: number, size: number) =>
  request(
    'TRANSACTIONS',
    api.getTransactions,
    (json: any) => normalize(json, { endpoint: 'currentPageTransactions' }),
    page,
    size
  )

export const fetchTransaction = (id: string) =>
  request('TRANSACTION', api.getTransaction, (json: any) => normalize(json), id)
