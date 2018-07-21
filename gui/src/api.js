import 'isomorphic-unfetch'
import formatRequestURI from 'utils/formatRequestURI'
import { AuthenticationError } from './errors'
import { camelizeKeys } from 'humps'

const formatURI = (path, query = {}) => {
  return formatRequestURI(path, query, {
    hostname: global.location.hostname,
    port: process.env.CHAINLINK_PORT
  })
}

const get = (path, query) => (
  global.fetch(
    formatURI(path, query),
    {credentials: 'include'}
  )
    .then(response => response.json())
    .then((data) => camelizeKeys(data))
)

const post = (path, body) => {
  return global.fetch(
    formatURI(path),
    {
      method: 'POST',
      body: JSON.stringify(body),
      credentials: 'include',
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json'
      }
    }
  )
    .then(response => {
      if (response.status === 401) {
        throw new AuthenticationError(response.statusText)
      }

      return response.json()
    })
    .then((data) => camelizeKeys(data))
}

const destroy = (path) => {
  return global.fetch(
    formatURI(path),
    {
      method: 'DELETE',
      credentials: 'include',
      headers: { 'Accept': 'application/json' }
    }
  )
    .then(response => response.json())
    .then((data) => camelizeKeys(data))
}

export const getJobs = (page, size) => get('/v2/specs', {page: page, size: size})
}
export const getJobSpec = (id) => get(`/v2/specs/${id}`)

export const getJobSpecRuns = (id, page, size) => get(`/v2/specs/${id}/runs`, {page: page, size: size})

export const getJobSpecRun = (id) => get(`/v2/runs/${id}`)

export const getAccountBalance = () => get('/v2/account_balance')

export const getConfiguration = () => get('/v2/config')

export const getBridges = (page, size) => get('/v2/bridge_types', {page: page, size: size})

export const getBridgeSpec = (name) => get(`/v2/bridge_types/${name}`)

export const createSession = (data) => post(`/sessions`, data)

export const destroySession = () => destroy(`/sessions`)
