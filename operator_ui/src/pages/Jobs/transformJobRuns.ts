import { ApiResponse } from 'utils/json-api-client'
import { JobRunV2 } from 'core/store/models'
import { parseDot, Stratify } from 'utils/parseDot'
import { PipelineJobRun, PipelineTaskRun } from './sharedTypes'
import { getJobStatus } from './utils'

function getTaskStatus({
  taskRun: { dotId, finishedAt, error },
  stratify,
  taskRuns,
}: {
  taskRun: JobRunV2['taskRuns'][0]
  stratify: Stratify[]
  taskRuns: JobRunV2['taskRuns']
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

const addTaskStatus =
  (stratify: Stratify[]) =>
  (
    taskRun: JobRunV2['taskRuns'][0],
    _index: number,
    taskRuns: JobRunV2['taskRuns'],
  ): PipelineTaskRun => {
    return {
      ...taskRun,
      status: getTaskStatus({ taskRun, stratify, taskRuns }),
    }
  }

export const transformPipelineJobRun =
  (jobSpecId: string) =>
  (jobRun: ApiResponse<JobRunV2>['data']): PipelineJobRun => {
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
      status: getJobStatus(jobRun.attributes),
      taskRuns,
      type: 'Pipeline job run',
    }
  }
