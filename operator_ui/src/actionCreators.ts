import * as jsonapi from '@chainlink/json-api-client'
import * as presenters from 'core/store/presenters'
import normalize from 'json-api-normalizer'
import { Action, Dispatch } from 'redux'
import { ThunkAction } from 'redux-thunk'
import * as api from './api'
import { Sessions } from './api/sessions'
import { RunStatus } from './core/store/models'
import { AppState } from './reducers'
import {
  AuthActionType,
  NotifyActionType,
  RedirectAction,
  ResourceActionType,
  RouterActionType,
} from './reducers/actions'

export type GetNormalizedData<T extends AnyFunc> = ReturnType<
  T
> extends ThunkAction<any, any, any, UpsertAction<infer A>>
  ? A
  : never

type Errors =
  | jsonapi.AuthenticationError
  | jsonapi.BadRequestError
  | jsonapi.ServerError
  | jsonapi.UnknownResponseError

const createErrorAction = (error: Error, type: string) => ({
  type,
  error: error.stack,
})

const REDIRECT_TO_SIGNOUT_ACTION: RedirectAction = {
  type: RouterActionType.REDIRECT,
  to: '/signout',
}

const curryErrorHandler = (dispatch: Dispatch, type: string) => (
  error: Error,
) => {
  if (error instanceof jsonapi.AuthenticationError) {
    dispatch(REDIRECT_TO_SIGNOUT_ACTION)
  } else {
    dispatch(createErrorAction(error, type))
  }
}

export const notifySuccess = (component: React.ReactNode, props: object) => {
  return {
    type: NotifyActionType.NOTIFY_SUCCESS,
    component,
    props,
  }
}

export const notifyError = (component: React.ReactNode, error: Error) => ({
  type: NotifyActionType.NOTIFY_ERROR,
  component,
  error,
})

/**
 * The type of any function
 */
type AnyFunc = (...args: any[]) => any

/**
 * Get the return type of a function, and unbox any promises
 */
type UnboxApi<T extends AnyFunc> = T extends (...args: any[]) => infer U
  ? U extends Promise<infer V>
    ? V
    : U
  : never

/**
 * Extract the first parameter from a function signature
 */
type Parameter<T extends AnyFunc> = Parameters<T>[0]

const signInSuccessAction = (doc: UnboxApi<Sessions['createSession']>) => {
  return {
    type: AuthActionType.RECEIVE_SIGNIN_SUCCESS,
    authenticated: doc.data.attributes.authenticated,
  }
}

const signInFailAction = () => ({ type: AuthActionType.RECEIVE_SIGNIN_FAIL })

function sendSignIn(data: Parameter<Sessions['createSession']>) {
  return (dispatch: Dispatch) => {
    dispatch({ type: AuthActionType.REQUEST_SIGNIN })

    return api.sessions
      .createSession(data)
      .then(doc => dispatch(signInSuccessAction(doc)))
      .catch((error: Errors) => {
        if (error instanceof jsonapi.AuthenticationError) {
          dispatch(signInFailAction())
        } else {
          dispatch(
            createErrorAction(error, AuthActionType.RECEIVE_SIGNIN_ERROR),
          )
        }
      })
  }
}

export const receiveSignoutSuccess = () => ({
  type: AuthActionType.RECEIVE_SIGNOUT_SUCCESS,
  authenticated: false,
})

function sendSignOut() {
  return (dispatch: Dispatch) => {
    dispatch({ type: AuthActionType.REQUEST_SIGNOUT })
    return api.sessions
      .destroySession()
      .then(() => dispatch(receiveSignoutSuccess()))
      .catch(curryErrorHandler(dispatch, AuthActionType.RECEIVE_SIGNIN_ERROR))
  }
}

const RECEIVE_CREATE_SUCCESS_ACTION = {
  type: ResourceActionType.RECEIVE_CREATE_SUCCESS,
}

const receiveDeleteSuccess = (id: string) => ({
  type: ResourceActionType.RECEIVE_DELETE_SUCCESS,
  id,
})

const receiveUpdateSuccess = (response: Response) => ({
  type: ResourceActionType.RECEIVE_UPDATE_SUCCESS,
  response,
})

export const submitSignIn = (data: Parameter<Sessions['createSession']>) =>
  sendSignIn(data)
export const submitSignOut = () => sendSignOut()

export const createJobSpec = (
  data: Parameter<typeof api.v2.specs.createJobSpec>,
  successCallback: React.ReactNode,
  errorCallback: React.ReactNode,
) => {
  return (dispatch: Dispatch) => {
    dispatch({ type: ResourceActionType.REQUEST_CREATE })

    return api.v2.specs
      .createJobSpec(data)
      .then(doc => {
        dispatch(RECEIVE_CREATE_SUCCESS_ACTION)
        dispatch(notifySuccess(successCallback, doc))
      })
      .catch((error: Errors) => {
        curryErrorHandler(
          dispatch,
          ResourceActionType.RECEIVE_CREATE_ERROR,
        )(error)
        dispatch(notifyError(errorCallback, error))
      })
  }
}

export const deleteJobSpec = (
  id: string,
  successCallback: React.ReactNode,
  errorCallback: React.ReactNode,
) => {
  return (dispatch: Dispatch) => {
    dispatch({ type: ResourceActionType.REQUEST_DELETE })

    return api.v2.specs
      .destroyJobSpec(id)
      .then(doc => {
        dispatch(receiveDeleteSuccess(id))
        dispatch(notifySuccess(successCallback, doc))
      })
      .catch((error: Errors) => {
        curryErrorHandler(
          dispatch,
          ResourceActionType.RECEIVE_DELETE_ERROR,
        )(error)
        dispatch(notifyError(errorCallback, error))
      })
  }
}

export const createJobRun = (
  id: string,
  successCallback: React.ReactNode,
  errorCallback: React.ReactNode,
): ThunkAction<Promise<void>, AppState, void, Action<string>> => {
  return (dispatch: Dispatch) => {
    dispatch({ type: ResourceActionType.REQUEST_CREATE })

    return api.v2.runs
      .createJobSpecRun(id)
      .then(doc => {
        dispatch(RECEIVE_CREATE_SUCCESS_ACTION)
        dispatch(notifySuccess(successCallback, doc))
      })
      .catch((error: Errors) => {
        curryErrorHandler(
          dispatch,
          ResourceActionType.RECEIVE_CREATE_ERROR,
        )(error)
        dispatch(notifyError(errorCallback, error))
      })
  }
}

export const createBridge = (
  data: Parameter<typeof api.v2.bridgeTypes.createBridge>,
  successCallback: React.ReactNode,
  errorCallback: React.ReactNode,
) => {
  return (dispatch: Dispatch) => {
    dispatch({ type: ResourceActionType.REQUEST_CREATE })

    return api.v2.bridgeTypes
      .createBridge(data)

      .then((doc: any) => {
        dispatch(RECEIVE_CREATE_SUCCESS_ACTION)
        dispatch(notifySuccess(successCallback, doc.data))
      })
      .catch((error: Errors) => {
        curryErrorHandler(
          dispatch,
          ResourceActionType.RECEIVE_CREATE_ERROR,
        )(error)
        dispatch(notifyError(errorCallback, error))
      })
  }
}

export const updateBridge = (
  params: Parameter<typeof api.v2.bridgeTypes.updateBridge>,
  successCallback: React.ReactNode,
  errorCallback: React.ReactNode,
) => {
  return (dispatch: Dispatch) => {
    dispatch({ type: ResourceActionType.REQUEST_UPDATE })

    return api.v2.bridgeTypes
      .updateBridge(params)
      .then((doc: any) => {
        dispatch(receiveUpdateSuccess(doc.data))
        dispatch(notifySuccess(successCallback, doc.data))
      })
      .catch((error: Errors) => {
        curryErrorHandler(
          dispatch,
          ResourceActionType.RECEIVE_UPDATE_ERROR,
        )(error)
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
  if (error instanceof jsonapi.AuthenticationError) {
    dispatch(REDIRECT_TO_SIGNOUT_ACTION)
  } else {
    dispatch(notifyError(({ msg }: any) => msg, error))
  }
}

/**
 * Extract the inner type of a promise if any
 */
type UnboxPromise<T> = T extends Promise<infer U> ? U : T

/**
 * An action to be dispatched to the store which contains the normalized
 * data from an external resource. Meant to be upserted into the required reducer.
 *
 * @template TNormalizedData The type of the normalized data to be upserted
 */
interface UpsertAction<TNormalizedData> extends Action<string> {
  data: TNormalizedData
}

/**
 * The request function is a factory function for async action creators.
 * Initially, a function is returned that accepts the arguments needed to make a request to an external resource.
 * Once this function is invoked, a thunk action is dispatched which invokes the request to the external resource.
 * One of two actions will occur on resource resolution:
 *
 * 1. When the resource is resolved successfully, we normalize the returned dataset and upsert it in the store.
 * 2. When the resource is resolved unsuccessfully, we handle the action via the `handleError` function
 *
 * @template TNormalizedData The shape of the output returned by the `normalizeData` function
 * @template TApiArgs The argument array to be fed to the `requestData` function, will be inferred from `requestData` parameter
 * @template TApiResp The response of the `requestData` function, will be inferred from `requestData` parameter
 *
 * @param type The action type field to be dispatched
 * @param requestData A function that outputs the data to be normalized and dispatched
 * @param normalizeData A function that normalizes the data returned by the requester function to be dispatched into an upsert action
 */
function request<
  TNormalizedData,
  TApiArgs extends Array<any>,
  TApiResp extends Promise<any>
>(
  type: string, // CHECKME -- stricten this type when we can
  requestData: (...args: TApiArgs) => TApiResp,
  normalizeData: (dataToNormalize: UnboxPromise<TApiResp>) => TNormalizedData,
): (
  ...args: TApiArgs
) => ThunkAction<
  Promise<void>,
  AppState,
  void,
  UpsertAction<TNormalizedData> | Action<string>
> {
  return (...args: TApiArgs) => {
    return dispatch => {
      dispatch({ type: `REQUEST_${type}` })

      return requestData(...args)
        .then(json => {
          const data = normalizeData(json)
          dispatch({ type: `UPSERT_${type}`, data })
        })
        .catch(handleError(dispatch))
        .finally(() => dispatch({ type: `RESPONSE_${type}` }))
    }
  }
}

export const fetchAccountBalance = request(
  'ACCOUNT_BALANCE',
  api.v2.user.balances.getAccountBalances,
  json =>
    normalize<{
      accountBalances: presenters.AccountBalance[]
    }>(json),
)

export type NormalizedAccountBalance = GetNormalizedData<
  typeof fetchAccountBalance
>

export const fetchConfiguration = request(
  'CONFIGURATION',
  api.v2.config.getConfiguration,
  normalize,
)

export const fetchBridges = request(
  'BRIDGES',
  api.v2.bridgeTypes.getBridges,
  json => normalize(json, { endpoint: 'currentPageBridges' }),
)

export const fetchBridgeSpec = request(
  'BRIDGE',
  api.v2.bridgeTypes.getBridgeSpec,
  json => normalize(json),
)

export const fetchJobs = request('JOBS', api.v2.specs.getJobSpecs, json =>
  normalize(json, { endpoint: 'currentPageJobs' }),
)

export const fetchRecentlyCreatedJobs = request(
  'RECENTLY_CREATED_JOBS',
  api.v2.specs.getRecentJobSpecs,
  json => normalize(json, { endpoint: 'recentlyCreatedJobs' }),
)

export const fetchJob = request('JOB', api.v2.specs.getJobSpec, json =>
  normalize(json, { camelizeKeys: false }),
)

export const fetchJobRuns = request(
  'JOB_RUNS',
  api.v2.runs.getJobSpecRuns,
  json => normalize(json, { endpoint: 'currentPageJobRuns' }),
)

export const fetchRecentJobRuns = request(
  'RECENT_JOB_RUNS',
  api.v2.runs.getRecentJobRuns,
  json => normalize(json, { endpoint: 'recentJobRuns' }),
)

export const fetchJobRun = request('JOB_RUN', api.v2.runs.getJobSpecRun, json =>
  normalize(json, { camelizeKeys: false }),
)

export const deleteCompletedJobRuns = (updatedBefore: string) =>
  request(
    'DELETE_COMPLETED_JOB_RUNS',
    api.v2.bulkDeleteRuns.bulkDeleteJobRuns,
    normalize,
  )({ status: [RunStatus.COMPLETED], updatedBefore })

export const deleteErroredJobRuns = (updatedBefore: string) =>
  request(
    'DELETE_ERRORED_JOB_RUNS',
    api.v2.bulkDeleteRuns.bulkDeleteJobRuns,
    normalize,
  )({ status: [RunStatus.ERRORED], updatedBefore })

export const fetchTransactions = request(
  'TRANSACTIONS',
  api.v2.transactions.getTransactions,
  json => normalize(json, { endpoint: 'currentPageTransactions' }),
)

export const fetchTransaction = request(
  'TRANSACTION',
  api.v2.transactions.getTransaction,
  json => normalize(json),
)

export const deleteJobSpecError = (id: number) =>
  request(
    'DELETE_JOB_SPEC_ERROR',
    api.v2.jobSpecErrors.destroyJobSpecError,
    normalize,
  )(id)
