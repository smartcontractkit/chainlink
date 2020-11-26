import TOML from '@iarna/toml'
import { PipelineTaskError, RunStatus } from 'core/store/models'

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
export const isOcrJob = (jobSpecId: string): boolean =>
  !isNaN((jobSpecId as unknown) as number)
