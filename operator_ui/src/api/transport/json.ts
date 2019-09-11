import {
  AttributesObject,
  ErrorsObject,
  JsonApiResponse,
  ResourceObject,
} from 'json-api-normalizer'
import pathToRegexp from 'path-to-regexp'
import fetchWithTimeout from 'utils/fetchWithTimeout'
import {
  AuthenticationError,
  BadRequestError,
  ServerError,
  UnknownResponseError,
} from '../errors'
import * as http from './http'

/**
 * The parameters required for making a paginated request to the chainlink endpoints
 */
export interface PaginatedRequestParams {
  size: number
  page: number
}

/**
 * A  json-api response for a data object.
 *
 *
 * @template T T is the type of the `attributes` object contained in the `data` key of `ResourceObject`, or null.
 * If T is an array of attribute objects, then the api response should be an array of resource objects.
 * If T is a single attribute object, then the api response is a single resource object, otherwise null.
 */
export interface ApiResponse<T extends AttributesObject | null>
  extends JsonApiResponse<
    T extends Array<infer U>
      ? ResourceObject<U>[]
      : T extends AttributesObject
      ? ResourceObject<T>
      : null
  > {}

/**
 * A paginated json-api response for a data object.
 *
 * The only difference between `PaginatedApiResponse` and `ApiResponse` is that `PaginatedApiResponse` includes a `TMeta` field
 * which supports the number of objects being returned, and a `TLinks` field supporting the previous and next page links.
 *
 * @template T T is the type of the `attributes` object contained in the `data` key of `ResourceObject`, or null.
 * If T is an array of attribute objects, then the api response should be an array of resource objects.
 * If T is a single attribute object, then the api response is a single resource object, otherwise null.
 */
export interface PaginatedApiResponse<T extends AttributesObject>
  extends JsonApiResponse<
    T extends Array<infer U> ? ResourceObject<U>[] : ResourceObject<T>,
    ErrorsObject[],
    never,
    {
      count: number
    },
    {
      prev?: string
      next?: string
    }
  > {}

export const fetchResource = methodFactory(http.Method.GET)
export const createResource = methodFactory(http.Method.POST)
export const updateResource = methodFactory(http.Method.PATCH)
export const deleteResource = methodFactory(http.Method.DELETE)

function methodFactory(method: http.Method) {
  return function<Params, T, NamedPathParams extends object = object>(
    url: string,
  ): Method<Params, T, NamedPathParams> {
    const toPath = pathToRegexp.compile<NamedPathParams>(url)

    return (params, namedPathParams) => {
      // if required, compile our path with its named path parameters
      const path = namedPathParams ? toPath(namedPathParams) : url
      const uri = http.formatURI(
        path,
        method === http.Method.GET ? params : undefined, // add query string options if its a GET method
      )
      const options = http.getOptions(method)
      const fetch = fetchWithTimeout(uri, options(params))

      return fetch.then(v => parseResponse(v))
    }
  }
}

/**
 * A json-api method which describes a function which accepts the required parameters to make a request,
 * and the return value of the request.
 *
 * @param params The parameters to the json-api endpoint
 * @param namedPathParams The named path parameters to the json-api endpoint
 * @template TParams The parameters to the json-api end point, it comes in the form of
 * an object which will either be serialized to the body of the request if it is a `POST`, `PATCH`, `DELETE` HTTP request,
 * or will be serialized to the query string of the url of the request if it is a `GET` HTTP request.
 * @template T The model of the data to be returned by the endpoint.
 * @template TNamedPathParams An object which the key will match the name of the path parameter of which the value will replace.
 * For example, for the path of `/v2/transaction/:txHash`, the value of `TNamedPathParams` should be:
 * ```ts
 * interface PathParams {
 *  txHash: string
 * }
 * ```
 */
type Method<TParams, T, TNamedPathParams extends object = object> = (
  params?: Partial<TParams>,
  namedPathParams?: TNamedPathParams,
) => Promise<ResponseType<TParams, T>>

/**
 * Our json-api response type is either paginated or non-paginated depending on
 * if the supplied api parameters extend `PaginatedRequestParams`
 *
 * @template TParams The parameters to the json-api end point, it comes in the form of
 * an object which will either be serialized to the body of the request if it is a `POST`, `PATCH`, `DELETE` HTTP request,
 * or will be serialized to the query string of the url of the request if it is a `GET` HTTP request.
 * @template T The model of the data to be returned by the endpoint.
 */
type ResponseType<TParams, T> = TParams extends PaginatedRequestParams
  ? PaginatedApiResponse<T>
  : ApiResponse<T>

function parseResponse<T>(response: Response): Promise<T> {
  if (response.status === 204) {
    return new Promise((resolve, _reject) => resolve({} as T))
  } else if (response.status >= 200 && response.status < 300) {
    return response.json()
  } else if (response.status === 400) {
    return response.json().then((doc: Pick<JsonApiResponse, 'errors'>) => {
      throw new BadRequestError(doc)
    })
  } else if (response.status === 401) {
    throw new AuthenticationError(response)
  } else if (response.status >= 500) {
    throw new ServerError(response)
  } else {
    throw new UnknownResponseError(response)
  }
}
