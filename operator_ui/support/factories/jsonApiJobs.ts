import { ApiResponse } from 'utils/json-api-client'
import { ResourceObject } from 'json-api-normalizer'
import { JobSpecV2 } from 'core/store/models'
import {
  cronJobV2,
  directRequestJobV2,
  fluxMonitorJobV2,
  keeperJobV2,
  ocrJobSpecV2,
  webhookJobV2,
  vrfJobV2,
} from './jobSpecV2'

function getRandomInt(max: number) {
  return Math.floor(Math.random() * Math.floor(max))
}

export const jsonApiJobSpecsV2 = (
  jobs: ResourceObject<JobSpecV2>[] = [],
  count?: number,
) => {
  const rc = count || jobs.length

  return {
    data: jobs,
    meta: { count: rc },
  } as ApiResponse<JobSpecV2[]>
}

export const directRequestResource = (
  job: Partial<JobSpecV2['directRequestSpec'] & { id?: string; name?: string }>,
) => {
  const id = job.id || getRandomInt(1_000_000).toString()

  return {
    type: 'jobs',
    id,
    attributes: {
      ...directRequestJobV2(job),
      name: job.name,
    },
  } as ResourceObject<JobSpecV2>
}

export const ocrJobResource = (
  job: Partial<
    JobSpecV2['offChainReportingOracleSpec'] & { id?: string; name?: string }
  >,
) => {
  const id = job.id || getRandomInt(1_000_000).toString()

  return {
    type: 'jobs',
    id,
    attributes: {
      ...ocrJobSpecV2(job),
      name: job.name,
    },
  } as ResourceObject<JobSpecV2>
}

export const fluxMonitorJobResource = (
  job: Partial<JobSpecV2['fluxMonitorSpec'] & { id?: string; name?: string }>,
) => {
  const id = job.id || getRandomInt(1_000_000).toString()

  return {
    type: 'jobs',
    id,
    attributes: fluxMonitorJobV2(job, { name: job.name }),
  } as ResourceObject<JobSpecV2>
}

export const keeperJobResource = (
  job: Partial<JobSpecV2['keeperSpec'] & { id?: string; name?: string }>,
) => {
  const id = job.id || getRandomInt(1_000_000).toString()

  return {
    type: 'jobs',
    id,
    attributes: keeperJobV2(job, { name: job.name }),
  } as ResourceObject<JobSpecV2>
}

export const cronJobResource = (
  job: Partial<JobSpecV2['cronSpec'] & { id?: string; name?: string }>,
) => {
  const id = job.id || getRandomInt(1_000_000).toString()

  return {
    type: 'jobs',
    id,
    attributes: cronJobV2(job, { name: job.name }),
  } as ResourceObject<JobSpecV2>
}

export const webJobResource = (
  job: Partial<JobSpecV2['webhookSpec'] & { id?: string; name?: string }>,
) => {
  const id = job.id || getRandomInt(1_000_000).toString()

  return {
    type: 'jobs',
    id,
    attributes: webhookJobV2(job, { name: job.name }),
  } as ResourceObject<JobSpecV2>
}

export const vrfJobResource = (
  job: Partial<JobSpecV2['vrfSpec'] & { id?: string; name?: string }>,
) => {
  const id = job.id || getRandomInt(1_000_000).toString()

  return {
    type: 'jobs',
    id,
    attributes: vrfJobV2(job, { name: job.name }),
  } as ResourceObject<JobSpecV2>
}
