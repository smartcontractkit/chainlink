import TOML from '@iarna/toml'

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
  value: string
  format: JobSpecFormats
}): string {
  if (format === JobSpecFormats.JSON) {
    try {
      return JSON.stringify(JSON.parse(value), null, 4)
    } catch {
      return value || ''
    }
  }

  return value
}
