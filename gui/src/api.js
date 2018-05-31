import { parse as parseQueryString } from 'query-string'
import url from 'url'
import { camelizeKeys } from 'humps'
import 'isomorphic-unfetch'

const DEFAULT_CHAINLINK_PORT = 6688

const formatUrl = path => {
  const port = parseQueryString(global.location.search).port || process.env.CHAINLINK_PORT || DEFAULT_CHAINLINK_PORT

  return url.format({
    hostname: global.location.hostname,
    port: port,
    pathname: path
  })
}

const request = path => (
  global.fetch(
    formatUrl(path),
    {credentials: 'include'}
  )
    .then(response => response.json())
    .then((data) => camelizeKeys(data))
)

export const getJobs = () => request('/v2/specs')

export const getAccountBalance = () => request('/v2/account_balance')
