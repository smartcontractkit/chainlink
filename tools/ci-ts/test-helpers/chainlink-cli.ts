import { ethers } from 'ethers'
import execa from 'execa'
import { JobSpec, JobRun } from '../../../operator_ui/@types/operator_ui'

const API_CREDENTIALS_PATH = '/run/secrets/apicredentials'

export function login(): void {
  execa.sync('chainlink', ['admin', 'login', '--file', API_CREDENTIALS_PATH])
}

export function getJobs(): JobSpec[] {
  const { stdout } = execa.sync('chainlink', ['-j', 'jobs', 'list'])
  return JSON.parse(stdout) as JobSpec[]
}

export function getJobRuns(): JobRun[] {
  const { stdout } = execa.sync('chainlink', ['-j', 'runs', 'list'])
  return JSON.parse(stdout) as JobRun[]
}

export function createJob(jobSpec: string): JobSpec {
  const { stdout } = execa.sync('chainlink', ['-j', 'jobs', 'create', jobSpec])
  return JSON.parse(stdout) as JobSpec
}

export function archiveJob(jobId: string): void {
  execa.sync('chainlink', ['-j', 'jobs', 'archive', jobId])
}

/**
 * interface for the data describing the CL node's keys
 */
interface KeyInfo {
  address: string
  ethBalance: ethers.utils.BigNumber
  linkBalance: ethers.utils.BigNumber
}

export function getAdminInfo(): KeyInfo[] {
  const { stdout } = execa.sync('chainlink', ['-j', 'admin', 'info'])
  return JSON.parse(stdout) as KeyInfo[]
}
