import { Api } from 'utils/json-api-client'
import { BridgeTypes } from './bridgeTypes'
import { BulkDeleteRuns } from './bulkDeleteRuns'
import { Chains } from './chains'
import { CSAKeys } from './csaKeys'
import { Config } from './config'
import { Features } from './features'
import { FeedsManagers } from './feedsManagers'
import { Jobs } from './jobs'
import { JobProposals } from './jobProposals'
import { OcrKeys } from './ocrKeys'
import { P2PKeys } from './p2pKeys'
import { Runs } from './runs'
import { Transactions } from './transactions'
import { User } from './user'
import { LogConfig } from './logConfig'
import { Nodes } from './nodes'

export class V2 {
  constructor(private api: Api) {}

  public bridgeTypes = new BridgeTypes(this.api)
  public bulkDeleteRuns = new BulkDeleteRuns(this.api)
  public chains = new Chains(this.api)
  public csaKeys = new CSAKeys(this.api)
  public config = new Config(this.api)
  public features = new Features(this.api)
  public feedsManagers = new FeedsManagers(this.api)
  public logConfig = new LogConfig(this.api)
  public nodes = new Nodes(this.api)
  public jobs = new Jobs(this.api)
  public jobProposals = new JobProposals(this.api)
  public ocrKeys = new OcrKeys(this.api)
  public p2pKeys = new P2PKeys(this.api)
  public runs = new Runs(this.api)
  public transactions = new Transactions(this.api)
  public user = new User(this.api)
}
