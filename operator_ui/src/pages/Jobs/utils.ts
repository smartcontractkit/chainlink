import TOML from '@iarna/toml'
import { PipelineTaskError, RunStatus } from 'core/store/models'
import { TaskSpec } from 'core/store/models'
import { parseDot, Stratify } from './parseDot'
import { countBy as _countBy } from 'lodash'

export function isToml({ value }: { value: string }): boolean {
  try {
    if (value !== '' && TOML.parse(value)) {
      return true
    } else {
      return false
    }
  } catch {
    return false
  }
}

export function stringifyJobSpec({
  value,
}: {
  value: { [key: string]: any }
}): string {
  try {
    return TOML.stringify(value)
  } catch (e) {
    console.error(
      `Failed to stringify job spec with the following error: ${e.message}`,
    )
    return ''
  }
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
  if (errorsExist(errors)) {
    return RunStatus.ERRORED
  }
  return RunStatus.COMPLETED
}

export function getTaskList({
  value,
}: {
  value: string
}): {
  list: false | TaskSpec[] | Stratify[]
  error: string
} {
  let list: false | TaskSpec[] | Stratify[] = false
  let error = ''

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

  return {
    list,
    error,
  }
}

function errorsExist(errors: PipelineTaskError[]): boolean {
  return errors !== null && errors.length > 0 && errors[0] !== null
}
