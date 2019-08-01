import { Dispatch } from 'redux'
import { pascalCase } from 'change-case'
import normalize from 'json-api-normalizer'
import * as api from './api'
import {
  AuthenticationError,
  BadRequestError,
  ServerError,
  UnknownResponseError
} from './api/errors'

type Errors =
  | AuthenticationError
  | BadRequestError
  | ServerError
  | UnknownResponseError

const createAction = (type: string) => ({ type: type })

const createErrorAction = (error: Error, type: string) => ({
  type: type,
  error: error.stack
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

interface Match {
  url: string
  params: object
}

export const matchRoute = (match: Match) => {
  return {
    type: MATCH_ROUTE,
    match: match
  }
}

export const NOTIFY_SUCCESS = 'NOTIFY_SUCCESS'

export const notifySuccess = (component: React.ReactNode, props: object) => {
  return {
    type: NOTIFY_SUCCESS,
    component: component,
    props: props
  }
}

export const NOTIFY_ERROR = 'NOTIFY_ERROR'

export const notifyError = (component: React.ReactNode, error: Error) => ({
  type: NOTIFY_ERROR,
  component: component,
  error: error
})

export const REQUEST_SIGNIN = 'REQUEST_SIGNIN'
export const RECEIVE_SIGNIN_SUCCESS = 'RECEIVE_SIGNIN_SUCCESS'
export const RECEIVE_SIGNIN_FAIL = 'RECEIVE_SIGNIN_FAIL'
export const RECEIVE_SIGNIN_ERROR = 'RECEIVE_SIGNIN_ERROR'

interface SignInDocument {
  data: {
    attributes: {
      authenticated: boolean
    }
  }
}

const signInSuccessAction = (doc: api.Document) => {
  const signDoc = <SignInDocument>doc

  return {
    type: RECEIVE_SIGNIN_SUCCESS,
    authenticated: signDoc.data.attributes.authenticated
  }
}

const signInFailAction = () => ({ type: RECEIVE_SIGNIN_FAIL })

function sendSignIn(data: object) {
  return (dispatch: Dispatch) => {
    dispatch(createAction(REQUEST_SIGNIN))
    return api
      .createSession(data)
      .then(doc => dispatch(signInSuccessAction(doc)))
      .catch((error: Errors) => {
        if (error instanceof AuthenticationError) {
          dispatch(signInFailAction())
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

const receiveCreateSuccessAction = () => ({ type: RECEIVE_CREATE_SUCCESS })

export const REQUEST_DELETE = 'REQUEST_DELETE'
export const RECEIVE_DELETE_SUCCESS = 'RECEIVE_DELETE_SUCCESS'
export const RECEIVE_DELETE_ERROR = 'RECEIVE_DELETE_ERROR'

const receiveDeleteSuccess = (id: string) => ({
  type: RECEIVE_DELETE_SUCCESS,
  id: id
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
  data: object,
  successCallback: React.ReactNode,
  errorCallback: React.ReactNode
) => {
  return (dispatch: Dispatch) => {
    dispatch(createAction(REQUEST_CREATE))
    return api
      .createJobSpec(data)
      .then(doc => {
        dispatch(receiveCreateSuccessAction())
        dispatch(notifySuccess(successCallback, doc))
      })
      .catch((error: Errors) => {
        curryErrorHandler(dispatch, RECEIVE_CREATE_ERROR)(error)
        dispatch(notifyError(errorCallback, error))
      })
  }
}

export const deleteJobSpec = (
  id: string,
  successCallback: React.ReactNode,
  errorCallback: React.ReactNode
) => {
  return (dispatch: Dispatch) => {
    dispatch(createAction(REQUEST_DELETE))
    return api
      .destroyJobSpec(id)
      .then(doc => {
        dispatch(receiveDeleteSuccess(id))
        dispatch(notifySuccess(successCallback, doc))
      })
      .catch((error: Errors) => {
        curryErrorHandler(dispatch, RECEIVE_DELETE_ERROR)(error)
        dispatch(notifyError(errorCallback, error))
      })
  }
}

export const createJobRun = (
  id: string,
  successCallback: React.ReactNode,
  errorCallback: React.ReactNode
) => {
  return (dispatch: Dispatch) => {
    dispatch(createAction(REQUEST_CREATE))
    return api
      .createJobSpecRun(id)
      .then((doc: any) => {
        dispatch(receiveCreateSuccessAction())
        dispatch(notifySuccess(successCallback, doc))
      })
      .catch((error: Errors) => {
        curryErrorHandler(dispatch, RECEIVE_CREATE_ERROR)(error)
        dispatch(notifyError(errorCallback, error))
      })
  }
}

export const createBridge = (
  data: object,
  successCallback: React.ReactNode,
  errorCallback: React.ReactNode
) => {
  return (dispatch: Dispatch) => {
    dispatch(createAction(REQUEST_CREATE))
    return api
      .createBridge(data)
      .then((doc: any) => {
        dispatch(receiveCreateSuccessAction())
        dispatch(notifySuccess(successCallback, doc.data))
      })
      .catch((error: Errors) => {
        curryErrorHandler(dispatch, RECEIVE_CREATE_ERROR)(error)
        dispatch(notifyError(errorCallback, error))
      })
  }
}

export const updateBridge = (
  params: api.UpdateBridgeParams,
  successCallback: React.ReactNode,
  errorCallback: React.ReactNode
) => {
  return (dispatch: Dispatch) => {
    dispatch(createAction(REQUEST_UPDATE))
    return api
      .updateBridge(params)
      .then((doc: any) => {
        dispatch(receiveUpdateSuccess(doc.data))
        dispatch(notifySuccess(successCallback, doc.data))
      })
      .catch((error: Errors) => {
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
  request('ACCOUNT_BALANCE', api.getAccountBalance, (json: object) =>
    normalize(json)
  )

export const fetchConfiguration = () =>
  request('CONFIGURATION', api.getConfiguration, (json: object) =>
    normalize(json)
  )

export const fetchBridges = (page: number, size: number) =>
  request(
    'BRIDGES',
    api.getBridges,
    (json: object) => normalize(json, { endpoint: 'currentPageBridges' }),
    page,
    size
  )

export const fetchBridgeSpec = (name: string) =>
  request('BRIDGE', api.getBridgeSpec, (json: object) => normalize(json), name)

export const fetchJobs = (page: number, size: number) =>
  request(
    'JOBS',
    api.getJobs,
    (json: object) => normalize(json, { endpoint: 'currentPageJobs' }),
    page,
    size
  )

export const fetchRecentlyCreatedJobs = (size: number) =>
  request(
    'RECENTLY_CREATED_JOBS',
    api.getRecentlyCreatedJobs,
    (json: object) => normalize(json, { endpoint: 'recentlyCreatedJobs' }),
    size
  )

export const fetchJob = (id: string) =>
  request('JOB', api.getJobSpec, (json: object) => normalize(json, { camelizeKeys: false }), id)

export const fetchJobRuns = (opts: api.JobSpecRunsOpts) =>
  request(
    'JOB_RUNS',
    api.getJobSpecRuns,
    (json: object) => normalize(json, { endpoint: 'currentPageJobRuns' }),
    opts
  )

export const fetchRecentJobRuns = (size: number) =>
  request(
    'RECENT_JOB_RUNS',
    api.getRecentJobRuns,
    (json: object) => normalize(json, { endpoint: 'recentJobRuns' }),
    size
  )

export const fetchJobRun = (id: string) =>
  request('JOB_RUN', api.getJobSpecRun, (json: object) => normalize(json, { camelizeKeys: false }), id)

export const deleteCompletedJobRuns = (updatedBefore: object) =>
  request(
    'DELETE_COMPLETED_JOB_RUNS',
    api.bulkDeleteJobRuns,
    normalize,
    ['completed'],
    updatedBefore
  )

export const deleteErroredJobRuns = (updatedBefore: object) =>
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
    (json: object) => normalize(json, { endpoint: 'currentPageTransactions' }),
    page,
    size
  )

export const fetchTransaction = (id: string) =>
  request(
    'TRANSACTION',
    api.getTransaction,
    (json: object) => normalize(json),
    id
  )
