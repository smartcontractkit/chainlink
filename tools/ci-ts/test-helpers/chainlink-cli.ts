import { ethers } from 'ethers'
import { exec } from 'shelljs'
import { JobSpec, RunResult } from '../../../operator_ui/@types/operator_ui'

function runCommand(command: string): Promise<string> {
  return new Promise((res, rej) => {
    exec(command, { silent: true }, (code, stdout, stderr) => {
      code === 0 ? res(stdout) : rej(stderr)
    })
  })
}

export async function login() {
  return runCommand('chainlink admin login --file /run/secrets/apicredentials')
}

export async function getJobs(): Promise<JobSpec[]> {
  const result = await runCommand('chainlink -j jobs list')
  return JSON.parse(result) as JobSpec[]
}

export async function getRunResults(): Promise<RunResult[]> {
  const result = await runCommand('chainlink -j runs list')
  return JSON.parse(result) as RunResult[]
}

export async function createJob(jobSpec: string): Promise<JobSpec> {
  const result = await runCommand(`chainlink -j jobs create '${jobSpec}'`)
  return JSON.parse(result) as JobSpec
}

export async function archiveJob(jobId: string): Promise<boolean> {
  await runCommand(`chainlink -j jobs archive '${jobId}'`)
  return true
}

interface KeyInfo {
  address: string
  ethBalance: ethers.utils.BigNumber
  linkBalance: ethers.utils.BigNumber
}

export async function getAdminInfo(): Promise<KeyInfo[]> {
  const result = await runCommand('chainlink -j admin info')
  return JSON.parse(result) as KeyInfo[]
}
