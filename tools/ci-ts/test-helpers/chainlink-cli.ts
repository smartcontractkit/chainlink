import { ethers } from 'ethers'
import execa from 'execa'
import { JobSpec, JobRun } from '../../../operator_ui/@types/operator_ui'
import crypto from 'crypto'
import path from 'path'
import Dockerode from 'dockerode'

const API_CREDENTIALS_PATH = '/run/secrets/apicredentials'
const docker = new Dockerode()

/**
 * interface for the data describing the CL node's keys
 */
interface KeyInfo {
  address: string
  ethBalance: ethers.utils.BigNumber
  linkBalance: ethers.utils.BigNumber
}

// make the root directory unique to the the URL and deterministic
function rootDirFromURL(value: string): string {
  const dirName = crypto
    .createHash('sha256')
    .update(value, 'utf8')
    .digest('hex')
    .slice(0, 16)
  return path.join('~', dirName)
}

export default class ChainlinkClient {
  name: string
  chainlinkURL: string
  container: Dockerode.Container
  root: string

  constructor(name: string, chainlinkURL: string, containerName: string) {
    this.name = name
    this.chainlinkURL = chainlinkURL
    this.container = docker.getContainer(containerName)
    this.root = rootDirFromURL(chainlinkURL)
  }

  async pause() {
    const paused = (await this.state()).Paused
    if (!paused) await this.container.pause()
  }

  async unpause() {
    const paused = (await this.state()).Paused
    if (paused) await this.container.unpause()
  }

  async state(): Promise<Dockerode.ContainerInspectInfo['State']> {
    return await this.container.inspect().then(res => res.State)
  }

  execute(command: string): object {
    const commands = ['-j', ...command.split(' ')]
    const { stdout } = execa.sync('chainlink', commands, this.execOptions())
    return stdout ? JSON.parse(stdout) : null
  }

  login(): void {
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
