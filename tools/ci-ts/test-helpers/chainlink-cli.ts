import { ethers } from 'ethers'
import execa from 'execa'
import { JobSpec, JobRun } from '../../../operator_ui/@types/operator_ui'

const API_CREDENTIALS_PATH = '/run/secrets/apicredentials'

export async function login(): Promise<void> {
  await execa('chainlink', ['admin', 'login', '--file', API_CREDENTIALS_PATH])
}

export async function getJobs(): Promise<JobSpec[]> {
  const { stdout } = await execa('chainlink', ['-j', 'jobs', 'list'])
  return JSON.parse(stdout) as JobSpec[]
}

export async function getJobRuns(): Promise<JobRun[]> {
  const { stdout } = await execa('chainlink', ['-j', 'runs', 'list'])
  return JSON.parse(stdout) as JobRun[]
}

export async function createJob(jobSpec: string): Promise<JobSpec> {
  const { stdout } = await execa('chainlink', ['-j', 'jobs', 'create', jobSpec])
  return JSON.parse(stdout) as JobSpec
}

export async function archiveJob(jobId: string): Promise<void> {
  await execa('chainlink', ['-j', 'jobs', 'archive', jobId])
}

/**
 * interface for the data describing the CL node's keys
 */
interface KeyInfo {
  address: string
  ethBalance: ethers.utils.BigNumber
  linkBalance: ethers.utils.BigNumber
}

export async function getAdminInfo(): Promise<KeyInfo[]> {
  const { stdout } = await execa('chainlink', ['-j', 'admin', 'info'])
  return JSON.parse(stdout) as KeyInfo[]
}
