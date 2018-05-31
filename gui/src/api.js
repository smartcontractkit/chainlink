import { parse as parseQueryString } from 'query-string'
import url from 'url'
import { camelizeKeys } from 'humps'
import 'isomorphic-unfetch'

const DEFAULT_CHAINLINK_PORT = 6688

const formatUrl = (path, query = {}) => {
  const port = parseQueryString(global.location.search).port || process.env.CHAINLINK_PORT || DEFAULT_CHAINLINK_PORT

  return url.format({
    hostname: global.location.hostname,
    port: port,
    pathname: path,
    query: query
  })
}

const request = (path, query) => (
  global.fetch(
    formatUrl(path, query),
    {credentials: 'include'}
  )
    .then(response => response.json())
    .then((data) => camelizeKeys(data))
)

export const getJobs = (page, size) => request('/v2/specs', {page: page, size: size})

export const getJobSpec = (id) => request(`/v2/specs/${id}`)

export const getAccountBalance = () => request('/v2/account_balance')
