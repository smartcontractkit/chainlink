import { JobSpecError, JobRunV2 } from 'core/store/models'
import * as time from 'time'

export type BaseJob = {
  createdAt: string
  definition: string
  errors: JobSpecError[]
  id: string
  name: string | null
}

export type JobSpecType =
  | 'directrequest'
  | 'fluxmonitor'
  | 'offchainreporting'
  | 'keeper'
  | 'cron'
  | 'webhook'
  | 'vrf'

export type JobV2 = BaseJob & {
  dotDagSource: string
  type: 'v2'
  specType: JobSpecType
}

export type BaseJobRun = {
  createdAt: time.Time
  finishedAt: time.Time | null
  id: string
  jobId: string
}

export type PipelineJobRunStatus = 'in_progress' | 'errored' | 'completed'
export type PipelineTaskRunStatus =
  | 'in_progress'
  | 'errored'
  | 'completed'
  | 'not_run'

export type PipelineTaskRun = JobRunV2['taskRuns'][0] & {
  status: PipelineTaskRunStatus
}

export type PipelineJobRun = BaseJobRun & {
  outputs: null | (string | null)[]
  errors: null | (string | null)[]
  pipelineSpec: {
    dotDagSource: string
  }
  status: PipelineJobRunStatus
  taskRuns: PipelineTaskRun[]
  type: 'Pipeline job run'
}

export type JobData = {
  job?: JobV2
  envAttributesDefinition?: string
  recentRuns?: PipelineJobRun[]
  recentRunsCount: number
  externalJobID?: string
}
