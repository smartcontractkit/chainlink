import {
  AttributesObject,
  ErrorsObject,
  JsonApiResponse,
  ResourceObject
} from 'json-api-normalizer'
import pathToRegexp from 'path-to-regexp'
import fetchWithTimeout from 'utils/fetchWithTimeout'
import {
  AuthenticationError,
  BadRequestError,
  ServerError,
  UnknownResponseError
} from '../errors'
import * as http from './http'

export interface PaginatedRequestParams {
  size: number
  page: number
}

export interface ApiResponse<T extends AttributesObject | null>
  extends JsonApiResponse<
    T extends Array<infer U>
      ? ResourceObject<U>[]
      : T extends AttributesObject
      ? ResourceObject<T>
      : null
  > {}

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
    url: string
  ) {
    const toPath = pathToRegexp.compile<NamedPathParams>(url)

    const createdMethod: Method<Params, T, NamedPathParams> = (
      params,
      namedPathParams
    ) => {
      const path = namedPathParams ? toPath(namedPathParams) : url
      const uri = http.formatURI(
        path,
        method === http.Method.GET ? params : undefined // add query string options if its a GET method
      )
      const options = http.getOptions(method)
      console.log(uri)
      const fetch = fetchWithTimeout(uri, options(params))
      const response = fetch.then(v =>
        parseResponse<ResponseType<Params, T>>(v)
      )

      return response
    }

    return createdMethod
  }
}

type Method<Params, T, NamedPathParams extends object = object> = (
  params?: Partial<Params>,
  namedPathParams?: NamedPathParams
) => Promise<ResponseType<Params, T>>

type ResponseType<Params, T> = Params extends PaginatedRequestParams
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
