import TOML from '@iarna/toml'
import { PipelineTaskError, RunStatus } from 'core/store/models'
import { TaskSpec } from 'core/store/models'
import { parseDot, Stratify } from './parseDot'
import { countBy as _countBy } from 'lodash'

export enum JobSpecFormats {
  JSON = 'json',
  TOML = 'toml',
}

export type JobSpecFormat = keyof typeof JobSpecFormats

export function isJson({ value }: { value: string }): JobSpecFormats | false {
  try {
    if (JSON.parse(value)) {
      return JobSpecFormats.JSON
    } else {
      return false
    }
  } catch {
    return false
  }
}

export function isToml({ value }: { value: string }): JobSpecFormats | false {
  try {
    if (value !== '' && TOML.parse(value)) {
      return JobSpecFormats.TOML
    } else {
      return false
    }
  } catch {
    return false
  }
}

export function getJobSpecFormat({
  value,
}: {
  value: string
}): JobSpecFormats | false {
  return isJson({ value }) || isToml({ value }) || false
}

export function stringifyJobSpec({
  value,
  format,
}: {
  value: { [key: string]: any }
  format: JobSpecFormats
}): string {
  try {
    if (format === JobSpecFormats.JSON) {
      return JSON.stringify(value, null, 4)
    } else if (format === JobSpecFormats.TOML) {
      return TOML.stringify(value)
    }
  } catch (e) {
    console.error(
      `Failed to stringify ${format} job spec with the following error: ${e.message}`,
    )
    return ''
  }

  return ''
}

export function getOcrJobStatus({
  finishedAt,
  errors,
}: {
  finishedAt: string | null
  errors: PipelineTaskError[]
}) {
  if (finishedAt === null) {
    return RunStatus.IN_PROGRESS
  }
  if (errors[0] !== null) {
    return RunStatus.ERRORED
  }
  return RunStatus.COMPLETED
}

// `isNaN` actually accepts strings and we don't want to `parseInt` or `parseFloat`
//  as it doesn't have the behaviour we want.
export const isJobV2 = (jobSpecId: string): boolean =>
  !isNaN((jobSpecId as unknown) as number)

export function getTaskList({
  value,
}: {
  value: string
}): {
  list: false | TaskSpec[] | Stratify[]
  format: false | JobSpecFormats
  error: string
} {
  const format = getJobSpecFormat({ value })
  let list: false | TaskSpec[] | Stratify[] = false
  let error = ''

  if (format === JobSpecFormats.JSON) {
    const tasks = JSON.parse(value).tasks
    list = Array.isArray(tasks) && tasks.length ? tasks : false
  } else if (format === JobSpecFormats.TOML) {
    try {
      const observationSource = parseDot(
        `digraph {${TOML.parse(value).observationSource as string}}`,
      )
      list =
        observationSource &&
        observationSource.length &&
        observationSource.some((node) => !node.parentIds.length)
          ? observationSource
          : false
      if (list) {
        list.every((listItem) => {
          const obj = _countBy(listItem.parentIds)
          Object.entries(obj).every(([parentId, value]) => {
            if (value > 1) {
              error += `${parentId} has duplicate ${listItem.id} children`
              list = false
              return false
            }

            return true
          })

          return !error
        })
      }
    } catch {
      list = false
    }
  }

  return {
    list,
    format,
    error,
  }
}
