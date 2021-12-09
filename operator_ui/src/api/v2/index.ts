import { Api } from 'utils/json-api-client'
import { BulkDeleteRuns } from './bulkDeleteRuns'
import { Chains } from './chains'
import { Config } from './config'
import { Jobs } from './jobs'
import { Transactions } from './transactions'
import { LogConfig } from './logConfig'
import { Nodes } from './nodes'
import { WebAuthn } from './webauthn'

export class V2 {
  constructor(private api: Api) {}

  public bulkDeleteRuns = new BulkDeleteRuns(this.api)
  public chains = new Chains(this.api)
  public config = new Config(this.api)
  public logConfig = new LogConfig(this.api)
  public nodes = new Nodes(this.api)
  public jobs = new Jobs(this.api)
  public transactions = new Transactions(this.api)
  public webauthn = new WebAuthn(this.api)
}
