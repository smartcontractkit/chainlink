import 'isomorphic-unfetch'
import formatRequestURI from 'utils/formatRequestURI'
import {
  AuthenticationError,
  BadRequestError,
  ServerError,
  UnknownResponseError
} from './errors'
import serializeBridgeType from 'api/serializers/bridgeType'

const formatURI = (path, query = {}) => {
  return formatRequestURI(path, query, {
    hostname: global.location.hostname,
    port: process.env.CHAINLINK_PORT
  })
}

const parseResponse = response => {
  if (response.status === 204) {
    return {}
  } else if (response.status >= 200 && response.status < 300) {
    return response.json()
  } else if (response.status === 400) {
    return response.json().then(json => {
      throw new BadRequestError(json)
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
    .fetch(formatURI(path, query), { credentials: 'include' })
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

const patch = (path, body) => {
  return global
    .fetch(formatURI(path), {
      method: 'PATCH',
      body: JSON.stringify(body),
      credentials: 'include',
      headers: {
        Accept: 'application/json',
        'Content-Type': 'application/json'
      }
    })
    .then(parseResponse)
}

const destroy = (path, body) =>
  global
    .fetch(formatURI(path, body), {
      method: 'DELETE',
      body: JSON.stringify(body),
      credentials: 'include',
      headers: { Accept: 'application/json' }
    })
    .then(parseResponse)

export const getJobs = (page, size) =>
  get('/v2/specs', { sort: '-createdAt', page: page, size: size })

export const getRecentlyCreatedJobs = size =>
  get('/v2/specs', { size: size, sort: '-createdAt' })

export const getJobSpec = id => get(`/v2/specs/${id}`)

export const getRecentJobRuns = size =>
  get(`/v2/runs`, { sort: '-createdAt', size: size })

export const getJobSpecRuns = (id, page, size) =>
  get(`/v2/runs`, { jobSpecId: id, sort: '-createdAt', page: page, size: size })

export const getJobSpecRun = id => get(`/v2/runs/${id}`)

export const getAccountBalance = () => get('/v2/user/balances')

export const getConfiguration = () => get('/v2/config')

export const getBridges = (page, size) =>
  get('/v2/bridge_types', { page: page, size: size })

export const getBridgeSpec = name => get(`/v2/bridge_types/${name}`)

export const createSession = data => post(`/sessions`, data)

export const createBridge = data => {
  const normalizedData = serializeBridgeType(data)
  return post('/v2/bridge_types', normalizedData)
}

export const createJobSpec = data => post('/v2/specs', data)

export const createJobSpecRun = id => post(`/v2/specs/${id}/runs`)

export const updateBridge = data => {
  const normalizedData = serializeBridgeType(data)
  return patch(`/v2/bridge_types/${data.name}`, normalizedData)
}

export const destroySession = () => destroy(`/sessions`)

export const bulkDeleteJobRuns = (status, updatedBefore) =>
  destroy('/v2/bulk_delete_runs', {
    status: status,
    updatedBefore: updatedBefore
  })

export const getTransactions = (page, size) =>
  get(`/v2/transactions`, { page: page, size: size })

export const getTransaction = id => get(`/v2/transactions/${id}`)
