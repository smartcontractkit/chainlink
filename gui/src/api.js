import { parse as parseQueryString } from 'query-string'
import url from 'url'
import { camelizeKeys } from 'humps'
import 'isomorphic-unfetch'

const formatUrl = (path, query = {}) => {
  let options = {
    hostname: global.location.hostname,
    pathname: path,
    query: query
  }

  const port = parseQueryString(global.location.search).port || process.env.CHAINLINK_PORT
  if (port) {
    options['port'] = port
  }
  return url.format(options)
}

const request = (path, query) => (
  global.fetch(
    formatUrl(path, query),
    {credentials: 'same-origin'}
  )
    .then(response => response.json())
    .then((data) => camelizeKeys(data))
)

export const getJobs = (page, size) => request('/v2/specs', {page: page, size: size})

export const getJobSpec = (id) => request(`/v2/specs/${id}`)

export const getJobSpecRuns = (id, page, size) => request(`/v2/specs/${id}/runs`, {page: page, size: size})

export const getJobSpecRun = (id) => request(`/v2/runs/${id}`)

export const getAccountBalance = () => request('/v2/account_balance')

export const getConfiguration = () => request('/v2/config')

export const getBridges = (page, size) => request('/v2/bridge_types', {page: page, size: size})
