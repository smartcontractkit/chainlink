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
