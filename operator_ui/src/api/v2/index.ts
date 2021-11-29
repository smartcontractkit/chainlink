import { Api } from 'utils/json-api-client'
import { BulkDeleteRuns } from './bulkDeleteRuns'
import { Chains } from './chains'
import { Config } from './config'
import { Features } from './features'
import { Jobs } from './jobs'
import { OcrKeys } from './ocrKeys'
import { P2PKeys } from './p2pKeys'
import { Runs } from './runs'
import { Transactions } from './transactions'
import { User } from './user'
import { LogConfig } from './logConfig'
import { Nodes } from './nodes'
import { WebAuthn } from './webauthn'

export class V2 {
  constructor(private api: Api) {}

  public bulkDeleteRuns = new BulkDeleteRuns(this.api)
  public chains = new Chains(this.api)
  public config = new Config(this.api)
  public features = new Features(this.api)
  public logConfig = new LogConfig(this.api)
  public nodes = new Nodes(this.api)
  public jobs = new Jobs(this.api)
  public ocrKeys = new OcrKeys(this.api)
  public p2pKeys = new P2PKeys(this.api)
  public runs = new Runs(this.api)
  public transactions = new Transactions(this.api)
  public user = new User(this.api)
  public webauthn = new WebAuthn(this.api)
}
