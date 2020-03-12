import { ethers } from 'ethers'
import execa from 'execa'
import { JobSpec, JobRun } from '../../../operator_ui/@types/operator_ui'

const API_CREDENTIALS_PATH = '/run/secrets/apicredentials'

/**
 * interface for the data describing the CL node's keys
 */
interface KeyInfo {
  address: string
  ethBalance: ethers.utils.BigNumber
  linkBalance: ethers.utils.BigNumber
}

export default class ChainlinkClient {
  chainlinkURL: string | undefined

  constructor(chainlinkURL?: string) {
    this.chainlinkURL = chainlinkURL
  }

  connect(chainlinkURL: string) {
    return new ChainlinkClient(chainlinkURL)
  }

  execute(command: string): object {
    if (!this.chainlinkURL) {
      throw Error('no chainlink node URL set')
    }
    const { stdout } = execa.sync('chainlink', ['-j', ...command.split(' ')])
    return stdout ? JSON.parse(stdout) : null
  }

  login(): void {
    execa.sync('chainlink', ['admin', 'login', '--file', API_CREDENTIALS_PATH])
  }

  getJobs(): JobSpec[] {
    // const { stdout } = execa.sync('chainlink', ['-j', 'jobs', 'list'])
    // return JSON.parse(stdout) as JobSpec[]
    return this.execute('jobs list') as JobSpec[]
  }

  getJobRuns(): JobRun[] {
    // const { stdout } = execa.sync('chainlink', ['-j', 'runs', 'list'])
    // return JSON.parse(stdout) as JobRun[]
    return this.execute('runs list') as JobRun[]
  }

  createJob(jobSpec: string): JobSpec {
    // const { stdout } = execa.sync('chainlink', [
    //   '-j',
    //   'jobs',
    //   'create',
    //   jobSpec,
    // ])
    // return JSON.parse(stdout) as JobSpec
    return this.execute(`jobs create ${jobSpec}`) as JobSpec
  }

  archiveJob(jobId: string): void {
    // execa.sync('chainlink', ['-j', 'jobs', 'archive', jobId])
    this.execute(`jobs archive ${jobId}`)
  }

  getAdminInfo(): KeyInfo[] {
    // const { stdout } = execa.sync('chainlink', ['-j', 'admin', 'info'])
    // return JSON.parse(stdout) as KeyInfo[]
    return this.execute('admin info') as KeyInfo[]
  }
}
