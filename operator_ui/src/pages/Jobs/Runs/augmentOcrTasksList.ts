import { PipelineTaskError, PipelineTaskOutput } from 'core/store/models'
import { parseDot, Stratify } from '../parseDot'
import { PipelineJobRun } from '../sharedTypes'

type AugmentedStratify = Stratify & {
  attributes: {
    error: PipelineTaskError
    output: PipelineTaskOutput
    status: PipelineJobRun['taskRuns'][0]['status']
    [key: string]: any
  }
}

function assignAttributes(stratify: Stratify): AugmentedStratify {
  if (stratify.attributes === undefined) {
    stratify.attributes = {}
  }

  return stratify as AugmentedStratify
}

export function augmentOcrTasksList({ jobRun }: { jobRun: PipelineJobRun }) {
  const graph = parseDot(`digraph {${jobRun.pipelineSpec.dotDagSource}}`)

  return graph.map((stratifyNode) => {
    const stratifyNodeCopy = assignAttributes(
      JSON.parse(JSON.stringify(stratifyNode)),
    )

    const taskRun = jobRun.taskRuns.find(
      ({ dotId }) => dotId === stratifyNodeCopy.id,
    )

    stratifyNodeCopy.attributes = {
      ...stratifyNodeCopy.attributes,
      error: taskRun?.error,
      output: taskRun?.output,
      status: taskRun?.status || 'not_run',
    }

    return stratifyNodeCopy
  })
}
