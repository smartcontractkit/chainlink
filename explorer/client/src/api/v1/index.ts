import { Api } from '@chainlink/json-api-client'
import { Auth } from './admin/auth'
import { Operators } from './admin/operators'
import { Heads } from './admin/heads'
import { JobRuns } from './jobRuns'
import { Config } from './config'

const api = new Api({
  base: process.env.REACT_APP_EXPLORER_BASEURL,
})

const adminAuth = new Auth(api)
const adminOperators = new Operators(api)
const adminHeads = new Heads(api)
const jobRuns = new JobRuns(api)
const config = new Config(api)

export { adminAuth, adminOperators, adminHeads, jobRuns, config }
