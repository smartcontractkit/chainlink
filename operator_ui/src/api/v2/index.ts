import { Api } from 'utils/json-api-client'
import { BulkDeleteRuns } from './bulkDeleteRuns'
import { Chains } from './chains'
import { Jobs } from './jobs'
import { LogConfig } from './logConfig'
import { Nodes } from './nodes'
import { WebAuthn } from './webauthn'

export class V2 {
  constructor(private api: Api) {}

  public bulkDeleteRuns = new BulkDeleteRuns(this.api)
  public chains = new Chains(this.api)
  public logConfig = new LogConfig(this.api)
  public nodes = new Nodes(this.api)
  public jobs = new Jobs(this.api)
  public webauthn = new WebAuthn(this.api)
}
