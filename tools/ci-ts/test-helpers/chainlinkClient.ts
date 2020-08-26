import { ethers } from 'ethers'
import execa from 'execa'
// TODO replace import with CL @types package
// https://www.pivotaltracker.com/story/show/171715396
import { JobSpec, JobRun } from '../../../operator_ui/@types/operator_ui'
import path from 'path'
import Dockerode from 'dockerode'

const docker = new Dockerode()

/**
 * interface for the data describing the CL node's keys
 */
interface KeyInfo {
  address: string
  ethBalance: ethers.utils.BigNumber
  linkBalance: ethers.utils.BigNumber
}

export default class ChainlinkClient {
  name: string
  clientURL: string
  container: Dockerode.Container
  rootDir: string

  private API_CREDENTIALS_PATH = '/run/secrets/apicredentials'

  constructor(name: string, clientURL: string, containerName: string) {
    this.name = name
    this.clientURL = clientURL
    this.container = docker.getContainer(containerName)
    this.rootDir = path.join('~', name)
  }

  /**
   * pauses the docker container running this CL node
   */
  public async pause() {
    await this.container.pause()
  }

  /**
   * unpauses the docker container running this CL node
   */
  public async unpause() {
    await this.container.unpause()
  }

  public login() {
    this.execute(`admin login --file ${this.API_CREDENTIALS_PATH}`)
  }

  public getJobs(): JobSpec[] {
    return this.execute('jobs list') as JobSpec[]
  }

  public getJobRuns(): JobRun[] {
    return this.execute('runs list') as JobRun[]
  }

  public createJob(jobSpec: string): JobSpec {
    return this.execute(`jobs create ${jobSpec}`) as JobSpec
  }

  public archiveJob(jobId: string): void {
    this.execute(`jobs archive ${jobId}`)
  }

  public getAdminInfo(): KeyInfo[] {
    return this.execute('admin info') as KeyInfo[]
  }

  /**
   * executes chainlink client commands within the docker image
   * @param command the command to pass to the chainlink CLI
   */
  private execute(command: string): object {
    const commands = ['-j', ...command.split(' ')]
    const { stdout } = execa.sync('chainlink', commands, this.execOptions())
    return stdout ? JSON.parse(stdout) : null
  }

  /**
   * options to pass to execa.sync command; useful for setting ENV variables,
   * which permits easily switching between different CL clients
   */
  private execOptions(): execa.SyncOptions {
    return {
      env: {
        CLIENT_NODE_URL: this.clientURL,
        ROOT: this.rootDir,
      },
    }
  }
}
