import { ApiResponse } from 'utils/json-api-client'
import { JobRun, OcrJobRun } from 'core/store/models'
import { parseDot, Stratify } from './parseDot'
import {
  DirectRequestJobRun,
  PipelineJobRun,
  PipelineTaskRun,
} from './sharedTypes'
import { getOcrJobStatus } from './utils'

function getTaskStatus({
  taskRun: { dotId, finishedAt, error },
  stratify,
  taskRuns,
}: {
  taskRun: OcrJobRun['taskRuns'][0]
  stratify: Stratify[]
  taskRuns: OcrJobRun['taskRuns']
}) {
  if (finishedAt === null) {
    return 'in_progress'
  }
  const currentNode = stratify.find((node) => node.id === dotId)

  let taskError = error

  if (currentNode) {
    currentNode.parentIds.forEach((id) => {
      const parentTaskRun = taskRuns.find((tr) => tr.dotId === id)

      if (parentTaskRun?.error !== null && parentTaskRun?.error === taskError) {
        taskError = 'not_run'
      }
    })
  }

  if (taskError === 'not_run') {
    return 'not_run'
  }

  if (taskError !== null) {
    return 'errored'
  }
  return 'completed'
}

const addTaskStatus = (stratify: Stratify[]) => (
  taskRun: OcrJobRun['taskRuns'][0],
  _index: number,
  taskRuns: OcrJobRun['taskRuns'],
): PipelineTaskRun => {
  return {
    ...taskRun,
    status: getTaskStatus({ taskRun, stratify, taskRuns }),
  }
}

export const transformPipelineJobRun = (jobSpecId: string) => (
  jobRun: ApiResponse<OcrJobRun>['data'],
): PipelineJobRun => {
  const stratify = parseDot(
    `digraph {${jobRun.attributes.pipelineSpec.dotDagSource}}`,
  )
  let taskRuns: PipelineTaskRun[] = []
  if (jobRun.attributes.taskRuns != null) {
    taskRuns = jobRun.attributes.taskRuns.map(addTaskStatus(stratify))
  }
  return {
    ...jobRun.attributes,
    id: jobRun.id,
    jobId: jobSpecId,
    status: getOcrJobStatus(jobRun.attributes),
    taskRuns,
    type: 'Pipeline job run',
  }
}

export const transformDirectRequestJobRun = (jobSpecId: string) => (
  jobRun: ApiResponse<JobRun>['data'],
): DirectRequestJobRun => ({
  ...jobRun.attributes,
  id: jobRun.id,
  jobId: jobSpecId,
  type: 'Direct request job run',
})
