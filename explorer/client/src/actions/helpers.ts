import * as jsonapi from '@chainlink/json-api-client'
import { Action } from 'redux'
import { ThunkAction } from 'redux-thunk'
import { AppState } from '../reducers'
import { FetchAdminSignoutSucceededAction } from '../reducers/actions'

/**
 * Extract the inner type of a promise if any
 */
export type UnboxPromise<T> = T extends Promise<infer U> ? U : T

/**
 * An action to be dispatched to the store which contains the normalized
 * data from an external resource. Meant to be upserted into the required reducer.
 *
 * @template TNormalizedData The type of the normalized data to be upserted
 */
export interface FetchAction<TNormalizedData> extends Action<string> {
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
export function request<
  TNormalizedData,
  TApiArgs extends Array<any>,
  TApiResp extends Promise<any>
>(
  type: string,
  requestData: (...args: TApiArgs) => TApiResp,
  normalizeData: (dataToNormalize: UnboxPromise<TApiResp>) => TNormalizedData,
): (
  ...args: TApiArgs
) => ThunkAction<
  Promise<void>,
  AppState,
  void,
  FetchAction<TNormalizedData> | Action<string>
> {
  return (...args: TApiArgs) => {
    return dispatch => {
      dispatch({ type: `FETCH_${type}_BEGIN` })

      return requestData(...args)
        .then(json => {
          const data = normalizeData(json)
          dispatch({ type: `FETCH_${type}_SUCCEEDED`, data })
        })
        .catch(e => {
          dispatch({ type: `FETCH_${type}_ERROR`, errors: e.errors })

          if (e instanceof jsonapi.AuthenticationError) {
            const fetchAdminSignoutSucceededAction: FetchAdminSignoutSucceededAction = {
              type: 'FETCH_ADMIN_SIGNOUT_SUCCEEDED',
            }
            dispatch(fetchAdminSignoutSucceededAction)
          }
        })
    }
  }
}
