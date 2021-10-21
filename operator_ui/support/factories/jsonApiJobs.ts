import { ApiResponse } from 'utils/json-api-client'
import { ResourceObject } from 'json-api-normalizer'
import { Job, JobSpecError } from 'core/store/models'
import {
  cronJobV2,
  directRequestJobV2,
  fluxMonitorJobV2,
  keeperJobV2,
  ocrJob,
  webhookJobV2,
  vrfJobV2,
} from './jobV2'

function getRandomInt(max: number) {
  return Math.floor(Math.random() * Math.floor(max))
}

export const jsonApiJob = (job: ResourceObject<Job>) => {
  return {
    data: job,
  } as ApiResponse<Job>
}

export const jsonApiJobSpecsV2 = (
  jobs: ResourceObject<Job>[] = [],
  count?: number,
) => {
  const rc = count || jobs.length

  return {
    data: jobs,
    meta: { count: rc },
  } as ApiResponse<Job[]>
}

export const directRequestResource = (
  job: Partial<Job['directRequestSpec'] & { id?: string; name?: string }>,
) => {
  const id = job.id || getRandomInt(1_000_000).toString()

  return {
    type: 'jobs',
    id,
    attributes: {
      ...directRequestJobV2(job),
      name: job.name,
    },
  } as ResourceObject<Job>
}

export const ocrJobResource = (
  job: Partial<
    Job['offChainReportingOracleSpec'] & { id?: string; name?: string }
  >,
) => {
  const id = job.id || getRandomInt(1_000_000).toString()

  return {
    type: 'jobs',
    id,
    attributes: {
      ...ocrJob(job),
      name: job.name,
    },
  } as ResourceObject<Job>
}

export const fluxMonitorJobResource = (
  job: Partial<
    Job['fluxMonitorSpec'] & {
      id?: string
      name?: string
      errors: JobSpecError[]
    }
  >,
) => {
  const id = job.id || getRandomInt(1_000_000).toString()

  return {
    type: 'jobs',
    id,
    attributes: fluxMonitorJobV2(job, { name: job.name, errors: job.errors }),
  } as ResourceObject<Job>
}

export const keeperJobResource = (
  job: Partial<Job['keeperSpec'] & { id?: string; name?: string }>,
) => {
  const id = job.id || getRandomInt(1_000_000).toString()

  return {
    type: 'jobs',
    id,
    attributes: keeperJobV2(job, { name: job.name }),
  } as ResourceObject<Job>
}

export const cronJobResource = (
  job: Partial<Job['cronSpec'] & { id?: string; name?: string }>,
) => {
  const id = job.id || getRandomInt(1_000_000).toString()

  return {
    type: 'jobs',
    id,
    attributes: cronJobV2(job, { name: job.name }),
  } as ResourceObject<Job>
}

export const webJobResource = (
  job: Partial<Job['webhookSpec'] & { id?: string; name?: string }>,
) => {
  const id = job.id || getRandomInt(1_000_000).toString()

  return {
    type: 'jobs',
    id,
    attributes: webhookJobV2(job, { name: job.name }),
  } as ResourceObject<Job>
}

export const vrfJobResource = (
  job: Partial<Job['vrfSpec'] & { id?: string; name?: string }>,
) => {
  const id = job.id || getRandomInt(1_000_000).toString()

  return {
    type: 'jobs',
    id,
    attributes: vrfJobV2(job, { name: job.name }),
  } as ResourceObject<Job>
}
