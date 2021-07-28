import { Api } from 'utils/json-api-client'
import { BridgeTypes } from './bridgeTypes'
import { BulkDeleteRuns } from './bulkDeleteRuns'
import { CSAKeys } from './csaKeys'
import { Config } from './config'
import { FeedsManagers } from './feedsManagers'
import { Runs } from './runs'
import { Specs } from './specs'
import { JobSpecErrors } from './jobSpecErrors'
import { Transactions } from './transactions'
import { User } from './user'
import { OcrKeys } from './ocrKeys'
import { P2PKeys } from './p2pKeys'
import { OcrRuns } from './ocrRuns'
import { Jobs } from './jobs'
import { JobProposals } from './jobProposals'
import { LogConfig } from './logConfig'

export class V2 {
  constructor(private api: Api) {}

  public bridgeTypes = new BridgeTypes(this.api)
  public bulkDeleteRuns = new BulkDeleteRuns(this.api)
  public config = new Config(this.api)
  public feedsManagers = new FeedsManagers(this.api)
  public runs = new Runs(this.api)
  public specs = new Specs(this.api)
  public jobProposals = new JobProposals(this.api)
  public jobSpecErrors = new JobSpecErrors(this.api)
  public transactions = new Transactions(this.api)
  public user = new User(this.api)
  public ocrKeys = new OcrKeys(this.api)
  public p2pKeys = new P2PKeys(this.api)
  public jobs = new Jobs(this.api)
  public ocrRuns = new OcrRuns(this.api)
  public logConfig = new LogConfig(this.api)
  public csaKeys = new CSAKeys(this.api)
}
