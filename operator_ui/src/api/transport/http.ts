import 'isomorphic-unfetch'
import formatRequestURI from 'src/utils/formatRequestURI'

export enum Method {
  GET = 'GET',
  POST = 'POST',
  PATCH = 'PATCH',
  DELETE = 'DELETE'
}

type FetchOptions = Parameters<typeof fetch>[1]

/**
 * Get the required request options for fetch
 * @param method The HTTP method to create options for
 */
export function getOptions(method: Method): (val?: any) => FetchOptions {
  if (method === Method.GET) {
    return () => ({
      method: 'GET',
      credentials: 'include'
    })
  }

  return CUDOptionFactory(method)
}

export function formatURI(path: string, query: object = {}) {
  return formatRequestURI(path, query, {
    hostname: location.hostname,
    port: process.env.CHAINLINK_PORT
  })
}

/**
 * Create, Update, Delete option factory
 * @param method The http method to create the option object for
 */
function CUDOptionFactory(method: Method): (body?: any) => FetchOptions {
  return (body?: any): FetchOptions => ({
    method,
    body: JSON.stringify(body),
    credentials: 'include',
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json'
    }
  })
}
