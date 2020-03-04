import { Api } from '@chainlink/json-api-client'
import { Auth } from './admin/auth'
import { Operators } from './admin/operators'
import { JobRuns } from './jobRuns'

const api = new Api({
  base: process.env.REACT_APP_EXPLORER_BASEURL,
})

const adminAuth = new Auth(api)
const adminOperators = new Operators(api)
const jobRuns = new JobRuns(api)

export { adminAuth, adminOperators, jobRuns }
