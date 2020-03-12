import { ethers } from 'ethers'
import execa from 'execa'
import { JobSpec, JobRun } from '../../../operator_ui/@types/operator_ui'
import crypto from 'crypto'
import path from 'path'

const API_CREDENTIALS_PATH = '/run/secrets/apicredentials'

/**
 * interface for the data describing the CL node's keys
 */
interface KeyInfo {
  address: string
  ethBalance: ethers.utils.BigNumber
  linkBalance: ethers.utils.BigNumber
}

function hashString(x: string): string {
  return crypto
    .createHash('sha256')
    .update(x, 'utf8')
    .digest('hex')
    .slice(0, 16)
}

export default class ChainlinkClient {
  chainlinkURL: string | undefined
  root: string | undefined

  constructor(chainlinkURL?: string) {
    if (chainlinkURL) {
      this.chainlinkURL = chainlinkURL
      // make the root directory unique to the the URL and deterministic
      this.root = path.join('~', hashString(chainlinkURL))
    }
  }

  connect(chainlinkURL: string) {
    return new ChainlinkClient(chainlinkURL)
  }

  execute(command: string): object {
    if (!this.chainlinkURL) {
      throw Error('no chainlink node URL set')
    }
    const commands = ['-j', ...command.split(' ')]
    const { stdout } = execa.sync('chainlink', commands, this.execOptions())
    return stdout ? JSON.parse(stdout) : null
  }

  login(): void {
    // execa.sync('chainlink', ['admin', 'login', '--file', API_CREDENTIALS_PATH])
    this.execute(`admin login --file ${API_CREDENTIALS_PATH}`)
  }

  getJobs(): JobSpec[] {
    return this.execute('jobs list') as JobSpec[]
  }

  getJobRuns(): JobRun[] {
    return this.execute('runs list') as JobRun[]
  }

  createJob(jobSpec: string): JobSpec {
    return this.execute(`jobs create ${jobSpec}`) as JobSpec
  }

  archiveJob(jobId: string): void {
    this.execute(`jobs archive ${jobId}`)
  }

  getAdminInfo(): KeyInfo[] {
    return this.execute('admin info') as KeyInfo[]
  }

  private execOptions(): execa.SyncOptions {
    return {
      env: {
        CLIENT_NODE_URL: this.chainlinkURL,
        ROOT: this.root,
      },
    }
  }
}
