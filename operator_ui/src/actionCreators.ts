import * as jsonapi from 'utils/json-api-client'
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
  ResourceActionType,
} from './reducers/actions'

export type GetNormalizedData<T extends AnyFunc> =
  ReturnType<T> extends ThunkAction<any, any, any, UpsertAction<infer A>>
    ? A
    : never

type Errors =
  | jsonapi.AuthenticationError
  | jsonapi.BadRequestError
  | jsonapi.ServerError
  | jsonapi.UnknownResponseError
  | jsonapi.ConflictError

type RestAction = 'UPSERT' | 'DELETE'

const createErrorAction = (error: Error, type: string) => ({
  type,
  error: error.stack,
})

const curryErrorHandler =
  (dispatch: Dispatch, type: string) => (error: Error) => {
    if (error instanceof jsonapi.AuthenticationError) {
      sendSignOut(dispatch)
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

export const notifySuccessMsg = (msg: string) => ({
  type: NotifyActionType.NOTIFY_SUCCESS_MSG,
  msg,
})

export const notifyError = (component: React.ReactNode, error: Error) => ({
  type: NotifyActionType.NOTIFY_ERROR,
  component,
  error,
})

export const notifyErrorMsg = (msg: string) => ({
  type: NotifyActionType.NOTIFY_ERROR_MSG,
  msg,
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
      .then((doc) => dispatch(signInSuccessAction(doc)))
      .catch((error: Errors) => {
        if (error instanceof jsonapi.AuthenticationError) {
          // Read the response to see if we're hitting a required MFA 401
          try {
            if (error.errors.length == 0 || error.errors[0].detail === null) {
              dispatch(signInFailAction())
              return
            }
            const errorResponse = error.errors[0].detail
            // Our response is good and we need to complete our challenge.
            errorResponse
              .json()
              .then((challengeData: any) => {
                if (!challengeData) {
                  // Ensure the data structure we're expecting is present
                  dispatch(signInFailAction())
                  return
                }

                // Throws if navigator is unavailable or user cancels flow
                try {
                  const publicKey = JSON.parse(
                    challengeData['errors'][0]['detail'],
                  )

                  publicKey.publicKey.challenge = bufferDecode(
                    publicKey.publicKey.challenge,
                  )
                  publicKey.publicKey.allowCredentials.forEach(
                    (listItem: any) => {
                      listItem.id = bufferDecode(listItem.id)
                    },
                  )

                  if (navigator.credentials === undefined) {
                    alert(
                      'Could not access credential subsystem in the browser. Must be using HTTPS or localhost.',
                    )
                    dispatch(signInFailAction())
                    return
                  }

                  navigator.credentials
                    .get({
                      publicKey: publicKey.publicKey,
                    })
                    .then((assertion: Credential | null) => {
                      if (assertion === null) {
                        // This likely means the user did not follow through
                        // with the attestation/authentication
                        dispatch(signInFailAction())
                        return
                      }

                      const pkassertion = assertion as PublicKeyCredential
                      const response =
                        pkassertion.response as AuthenticatorAssertionResponse

                      const authData = response.authenticatorData
                      const clientDataJSON = response.clientDataJSON
                      const rawId = pkassertion.rawId
                      const sig = response.signature
                      const userHandle = response.userHandle

                      // Build our response assertion
                      const waData = JSON.stringify({
                        id: assertion.id,
                        rawId: bufferEncode(rawId),
                        type: assertion.type,
                        response: {
                          authenticatorData: bufferEncode(authData),
                          clientDataJSON: bufferEncode(clientDataJSON),
                          signature: bufferEncode(sig),
                          userHandle: bufferEncode(userHandle),
                        },
                      })

                      data.webauthndata = waData

                      // Retry login with this new attestation
                      return api.sessions
                        .createSession(data)
                        .then((doc) => dispatch(signInSuccessAction(doc)))
                        .catch((_error: Errors) => {
                          dispatch(signInFailAction())
                        })
                    })
                } catch {
                  dispatch(signInFailAction())
                }
              })
              .catch((_error: Errors) => {
                // The detail field was not parsable JSON
                dispatch(signInFailAction())
              })
          } catch {
            // There was no data in our 401 response. So this is just a bad password
            dispatch(signInFailAction())
          }
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

function sendSignOut(dispatch: Dispatch) {
  return api.sessions
    .destroySession()
    .then(() => dispatch(receiveSignoutSuccess()))
    .catch(curryErrorHandler(dispatch, AuthActionType.RECEIVE_SIGNIN_ERROR))
}

// Base64 to ArrayBuffer
function bufferDecode(value: any) {
  return Uint8Array.from(atob(value), (c) => c.charCodeAt(0))
}

// ArrayBuffer to URLBase64
function bufferEncode(value: ArrayBuffer | null) {
  if (value === null) {
    return ''
  }

  const uint8View = new Uint8Array(value)
  const ar = String.fromCharCode.apply(null, Array.from(uint8View))
  return btoa(ar).replace(/\+/g, '-').replace(/\//g, '_').replace(/=/g, '')
}

function completeKeyRegistration(response: any) {
  const credentialCreationOptions = response['data']['attributes']['settings']
  credentialCreationOptions.publicKey.challenge = bufferDecode(
    credentialCreationOptions.publicKey.challenge,
  )
  credentialCreationOptions.publicKey.user.id = bufferDecode(
    credentialCreationOptions.publicKey.user.id,
  )
  if (credentialCreationOptions.publicKey.excludeCredentials) {
    credentialCreationOptions.publicKey.excludeCredentials.forEach(
      (excludeCredential: any) => {
        excludeCredential.id = bufferDecode(excludeCredential.id)
      },
    )
  }

  return navigator.credentials.create({
    publicKey: credentialCreationOptions.publicKey,
  })
}

function sendBeginRegistration() {
  if (navigator.credentials === undefined) {
    alert(
      'Could not access credential subsystem in the browser. Must be using HTTPS or localhost.',
    )
    return
  }

  return api.v2.webauthn
    .beginKeyRegistration({})
    .then((response) =>
      completeKeyRegistration(response).then(
        (credential: Credential | null) => {
          if (credential === null) {
            alert(
              'Error, could not generate credential. User declined to enroll?',
            )
            return
          }

          const pkcredential = credential as PublicKeyCredential
          const response =
            pkcredential.response as AuthenticatorAttestationResponse

          const credentialStr = {
            id: credential.id,
            rawId: bufferEncode(pkcredential.rawId),
            type: credential.type,
            response: {
              attestationObject: bufferEncode(response.attestationObject),
              clientDataJSON: bufferEncode(response.clientDataJSON),
            },
          }
          return api.v2.webauthn.finishKeyRegistration(credentialStr)
        },
      ),
    )
    .catch((err) => {
      alert(
        'Key registration error, ensure MFA_RPID and MFA_RPORIGIN environment variables are set.\n' +
          err,
      )
    })
}

const RECEIVE_CREATE_SUCCESS_ACTION = {
  type: ResourceActionType.RECEIVE_CREATE_SUCCESS,
}

const receiveDeleteSuccess = (id: string) => ({
  type: ResourceActionType.RECEIVE_DELETE_SUCCESS,
  id,
})

export const submitSignIn = (data: Parameter<Sessions['createSession']>) =>
  sendSignIn(data)

export const submitSignOut = () => sendSignOut

export const beginRegistration = () => sendBeginRegistration()

export const deleteChain = (
  id: string,
  successCallback: React.ReactNode,
  errorCallback: React.ReactNode,
) => {
  return (dispatch: Dispatch) => {
    dispatch({ type: ResourceActionType.REQUEST_DELETE })

    const endpoint = api.v2.chains

    return endpoint
      .destroyChain(id)
      .then((doc) => {
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

export const createJobRunV2 = (
  id: string,
  pipelineInput: string,
  successCallback: React.ReactNode,
  errorCallback: React.ReactNode,
): ThunkAction<Promise<void>, AppState, void, Action<string>> => {
  return (dispatch: Dispatch) => {
    dispatch({ type: ResourceActionType.REQUEST_CREATE })

    return api.v2.jobs
      .createJobRunV2(id, pipelineInput)
      .then((doc) => {
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

// DEV NOTE:
// Above here is deprecated. Use the `request(...)` function below to wrap API
// calls in a counter, normalize JSON-API responses and create notifications.
//
// The calls above will be converted gradually.
const handleError = (dispatch: Dispatch) => (error: Error) => {
  if (error instanceof jsonapi.AuthenticationError) {
    dispatch(receiveSignoutSuccess())
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
 * @param prefix A prefix to apply to the request actions to signify the REST action
 * @param requestData A function that outputs the data to be normalized and dispatched
 * @param normalizeData A function that normalizes the data returned by the requester function to be dispatched into an upsert action
 */
function request<
  TNormalizedData,
  TApiArgs extends Array<any>,
  TApiResp extends Promise<any>,
>(
  type: string, // CHECKME -- stricten this type when we can
  prefix: RestAction,
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
  const requestType =
    prefix === 'UPSERT' ? `REQUEST_${type}` : `REQUEST_${prefix}_${type}`
  const responseType =
    prefix === 'UPSERT' ? `RESPONSE_${type}` : `RESPONSE_${prefix}_${type}`
  const successType = `${prefix}_${type}`

  return (...args: TApiArgs) => {
    return (dispatch) => {
      dispatch({ type: requestType })

      return requestData(...args)
        .then((json) => {
          const data = normalizeData(json)
          dispatch({ type: successType, data })
        })
        .catch(handleError(dispatch))
        .finally(() => dispatch({ type: responseType }))
    }
  }
}

/**
 * requestDelete calls the request action creator, specifying 'DELETE' as the action prefix
 *
 * @param type The action type field to be dispatched
 * @param requestData A function that outputs the data to be normalized and dispatched
 * @param normalizeData A function that normalizes the data returned by the requester function to be dispatched into an upsert action
 */
function requestDelete(
  type: Parameters<typeof request>[0],
  requestData: Parameters<typeof request>[2],
  normalizeData: Parameters<typeof request>[3],
): ReturnType<typeof request> {
  return request(type, 'DELETE', requestData, normalizeData)
}

export const deleteErroredJobRuns = (updatedBefore: string) =>
  requestDelete(
    'ERRORED_JOB_RUNS',
    api.v2.bulkDeleteRuns.bulkDeleteJobRuns,
    normalize,
  )({ status: [RunStatus.ERRORED], updatedBefore })
