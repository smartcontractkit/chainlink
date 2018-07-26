import formatRequestURI from 'utils/formatRequestURI'
import { camelizeKeys } from 'humps'
import 'isomorphic-unfetch'

const formatURI = (path, query = {}) => {
  return formatRequestURI(path, query, {
    hostname: global.location.hostname,
    port: process.env.CHAINLINK_PORT
  })
}

const request = (path, query) => (
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
    .then(response => response.json())
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

export const getJobs = (page, size) => request('/v2/specs', {page: page, size: size})

export const getJobSpec = (id) => request(`/v2/specs/${id}`)

export const getJobSpecRuns = (id, page, size) => request(`/v2/specs/${id}/runs`, {page: page, size: size})

export const getJobSpecRun = (id) => request(`/v2/runs/${id}`)

export const getAccountBalance = () => request('/v2/account_balance')

export const getConfiguration = () => request('/v2/config')

export const getBridges = (page, size) => request('/v2/bridge_types', {page: page, size: size})

export const getBridgeSpec = (name) => request(`/v2/bridge_types/${name}`)

export const postSessionRequest = (data) => post(`/sessions`, data)

export const destroySession = () => destroy(`/sessions`)
