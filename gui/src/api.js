import { parse as parseQueryString } from 'query-string'
import url from 'url'
import { camelizeKeys } from 'humps'
import 'isomorphic-unfetch'

const DEFAULT_CHAINLINK_PORT = 6688

const formatUrl = (path) => {
  const port = parseQueryString(global.location.search).port || process.env.CHAINLINK_PORT || DEFAULT_CHAINLINK_PORT

  return url.format({
    hostname: global.location.hostname,
    port: port,
    pathname: path
  })
}

export const getJobs = () => {
  const requestUrl = formatUrl('/v2/specs')

  return global.fetch(requestUrl, {credentials: 'include'})
    .then(response => response.json())
    .then((data) => camelizeKeys(data))
}

export const getAccountBalance = () => {
  const requestUrl = formatUrl('/v2/account_balance')

  return global.fetch(requestUrl, {credentials: 'include'})
    .then(response => response.json())
    .then((data) => camelizeKeys(data))
}
