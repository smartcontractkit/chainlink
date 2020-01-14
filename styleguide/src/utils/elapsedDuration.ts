const MS_IN_SECOND = 1000
const SECONDS_IN_MINUTE = 60
const SECONDS_IN_HOUR = 3600

export type DateValue = string | number
export type FinishedAt = DateValue | null

export function elapsedDuration(
  createdAt: DateValue,
  finishedAt: FinishedAt,
): string {
  if (createdAt === '' && finishedAt === '') {
    return ''
  }

  const es = elapsedSeconds(new Date(createdAt), endAt(finishedAt))
  const hours = Math.floor(es / SECONDS_IN_HOUR)
  const minutes = Math.floor((es % SECONDS_IN_HOUR) / SECONDS_IN_MINUTE)
  const seconds = Math.ceil((es % SECONDS_IN_HOUR) % SECONDS_IN_MINUTE)

  return format(hours, minutes, seconds)
}

function endAt(finishedAt: FinishedAt): Date {
  if (finishedAt === null) {
    return new Date(Date.now())
  }

  return new Date(finishedAt)
}

function elapsedSeconds(from: Date, to: Date): number {
  return to.getTime() / MS_IN_SECOND - from.getTime() / MS_IN_SECOND
}

function format(hours: number, minutes: number, seconds: number): string {
  return `${hours ? `${hours}h` : ''}${minutes ? `${minutes}m` : ''}${seconds}s`
}
