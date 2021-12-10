import TOML from '@iarna/toml'
import { RunStatus } from 'core/store/models'
import { TaskSpec } from 'core/store/models'
import { parseDot, Stratify } from 'utils/parseDot'
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

// Temporary function until we can come up with a better design
export function getJobStatusGQL({
  finishedAt,
  errors,
}: {
  finishedAt: string | null
  errors: ReadonlyArray<string>
}) {
  if (finishedAt === null) {
    return RunStatus.IN_PROGRESS
  }
  if (errors !== null && errors.length > 0 && errors[0] !== null) {
    return RunStatus.ERRORED
  }
  return RunStatus.COMPLETED
}

export function getTaskList({ value }: { value: string }): {
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
