import 'isomorphic-unfetch'

export enum Method {
  GET = 'GET',
  POST = 'POST',
  PATCH = 'PATCH',
  DELETE = 'DELETE',
}

type FetchOptions = Parameters<typeof fetch>[1]

/**
 * Get the required request options for fetch
 * @param method The HTTP method to create options for
 * @param raw If true, the body will not be JSON stringified
 */
export function getOptions(
  method: Method,
  raw?: boolean,
): (val?: any) => FetchOptions {
  if (method === Method.GET) {
    return () => ({
      method: 'GET',
      credentials: 'include',
    })
  }

  return CUDOptionFactory(method, raw)
}

/**
 * Create, Update, Delete option factory
 * @param method The http method to create the option object for
 * @param raw If true, the body will not be JSON stringified
 */
function CUDOptionFactory(
  method: Method,
  raw?: boolean,
): (body?: any) => FetchOptions {
  return (body?: any): FetchOptions => ({
    method,
    body: raw ? body : JSON.stringify(body),
    credentials: 'include',
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json',
    },
  })
}

export function createUrl(base: string, path: string, query?: object) {
  let u: URL
  try {
    u = new URL(path, base)
  } catch (e) {
    const origMsg = (e as TypeError).message
    const newMsg = `Error when creating url with path=${path}, base=${base}: ${origMsg}`
    const err = Error(newMsg)
    throw err
  }

  if (query) {
    Object.entries(query).forEach(([k, v]) => {
      if (typeof v === 'string') {
        return u.searchParams.append(k, v)
      }

      // serialize v as long as its not "null" or undefined
      if (v != undefined) {
        // eslint-disable-next-line @typescript-eslint/ban-types
        return u.searchParams.append(k, (v as Object).toString())
      }
    })
  }
  return u
}
