import normalize from 'json-api-normalizer'
import * as api from '../api/index'
import { request } from './helpers'

export const fetchJobRuns = request(
  'JOB_RUNS',
  api.v1.jobRuns.getJobRuns,
  json => normalize(json, { endpoint: 'currentPageJobRuns' }),
)

export const fetchJobRun = request('JOB_RUN', api.v1.jobRuns.getJobRun, json =>
  normalize(json, { endpoint: 'jobRun' }),
)
