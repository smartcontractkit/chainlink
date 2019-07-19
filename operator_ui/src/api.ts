import 'isomorphic-unfetch'
import {
  AuthenticationError,
  BadRequestError,
  DocumentWithErrors,
  ServerError,
  UnknownResponseError
} from './api/errors'
import serializeBridgeType from './api/serializers/bridgeType'
import fetchWithTimeout from './utils/fetchWithTimeout'
import formatRequestURI from './utils/formatRequestURI'

const formatURI = (path: string, query: object = {}) => {
  return formatRequestURI(path, query, {
    hostname: location.hostname,
    port: process.env.CHAINLINK_PORT
  })
}

interface EmptyDocument {}

export interface DocumentData {
  id: any
  attributes: { [k: string]: any }
}

export interface DocumentWithData {
  data: DocumentData | DocumentData[]
  links?: object
  meta?: object
}

export type Document = DocumentWithData | EmptyDocument

const parseResponse = (response: Response): Promise<Document> => {
  if (response.status === 204) {
    return new Promise((resolve, _reject) => resolve({} as EmptyDocument))
  } else if (response.status >= 200 && response.status < 300) {
    return response.json()
  } else if (response.status === 400) {
    return response.json().then((doc: DocumentWithErrors) => {
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

const get = (path, query) =>
  global
    .fetchWithTimeout(formatURI(path, query), { credentials: 'include' })
    .then(parseResponse)

const post = (path, body) => {
  return global
    .fetch(formatURI(path), {
      method: 'POST',
      body: JSON.stringify(body),
      credentials: 'include',
      headers: {
        Accept: 'application/json',
        'Content-Type': 'application/json'
      }
    })
    .then(parseResponse)
}

const patch = (path: string, body: object) => {
  return fetch(formatURI(path), {
    method: 'PATCH',
    body: JSON.stringify(body),
    credentials: 'include',
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json'
    }
  }).then(parseResponse)
}

const destroy = (path: string, body: object = {}) => {
  return fetch(formatURI(path, body), {
    method: 'DELETE',
    body: JSON.stringify(body),
    credentials: 'include',
    headers: { Accept: 'application/json' }
  }).then(parseResponse)
}

export const getJobs = (page: number, size: number) =>
  get('/v2/specs', { sort: '-createdAt', page: page, size: size })

export const getRecentlyCreatedJobs = (size: number) =>
  get('/v2/specs', { size: size, sort: '-createdAt' })

export const getJobSpec = (id: string) => get(`/v2/specs/${id}`)

export const getRecentJobRuns = (size: number) =>
  get(`/v2/runs`, { sort: '-createdAt', size: size })

export interface JobSpecRunsOpts {
  page: number
  size: number
  sort?: string
  jobSpecId?: string
}

export const getJobSpecRuns = ({
  jobSpecId,
  page,
  size,
  sort = '-createdAt'
}: JobSpecRunsOpts) => {
  const opts = { sort, page, size } as JobSpecRunsOpts
  if (jobSpecId) {
    opts.jobSpecId = jobSpecId
  }

  return get(`/v2/runs`, opts)
}

export const getJobSpecRun = (id: string) => get(`/v2/runs/${id}`)

export const getAccountBalance = () => get('/v2/user/balances')

export const getConfiguration = () => get('/v2/config')

export const getBridges = (page: number, size: number) =>
  get('/v2/bridge_types', { page: page, size: size })

export const getBridgeSpec = (name: string) => get(`/v2/bridge_types/${name}`)

export const createSession = (data: object) => post(`/sessions`, data)

export const createBridge = (data: object) => {
  const normalizedData = serializeBridgeType(data)
  return post('/v2/bridge_types', normalizedData)
}

export const createJobSpec = (data: object) => post('/v2/specs', data)

export const createJobSpecRun = (id: string) => post(`/v2/specs/${id}/runs`)

export interface UpdateBridgeParams {
  name: string
  url: string
  confirmations: number
  minimumContractPayment?: string
}

export const updateBridge = (params: UpdateBridgeParams) => {
  const normalizedData = serializeBridgeType(params)
  return patch(`/v2/bridge_types/${params.name}`, normalizedData)
}

export const destroyJobSpec = (id: string) => destroy(`/v2/specs/${id}`)

export const destroySession = () => destroy(`/sessions`)

export const bulkDeleteJobRuns = (status: string, updatedBefore: string) =>
  destroy('/v2/bulk_delete_runs', {
    status: status,
    updatedBefore: updatedBefore
  })

export const getTransactions = (page: number, size: number) =>
  get(`/v2/transactions`, { page: page, size: size })

export const getTransaction = (id: string) => get(`/v2/transactions/${id}`)
