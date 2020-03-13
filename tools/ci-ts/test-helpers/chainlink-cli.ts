import { ethers } from 'ethers'
import execa from 'execa'
import { JobSpec, JobRun } from '../../../operator_ui/@types/operator_ui'
import crypto from 'crypto'
import path from 'path'
import Docker from 'dockerode'

const API_CREDENTIALS_PATH = '/run/secrets/apicredentials'
const docker = new Docker()

/**
 * interface for the data describing the CL node's keys
 */
interface KeyInfo {
  address: string
  ethBalance: ethers.utils.BigNumber
  linkBalance: ethers.utils.BigNumber
}

function hashString(value: string): string {
  return crypto
    .createHash('sha256')
    .update(value, 'utf8')
    .digest('hex')
    .slice(0, 16)
}

export default class ChainlinkClient {
  name: string
  chainlinkURL: string
  container: Docker.Container
  root: string

  constructor(name: string, chainlinkURL: string, containerName: string) {
    this.name = name
    this.chainlinkURL = chainlinkURL
    this.container = docker.getContainer(containerName)
    // make the root directory unique to the the URL and deterministic
    this.root = path.join('~', hashString(chainlinkURL))
  }

  async pause() {
    await this.container.pause()
  }

  async unpause() {
    await this.container.unpause()
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
