const MS_IN_SECOND = 1000
const SECONDS_IN_MINUTE = 60
const SECONDS_IN_HOUR = 3600

export function elapsedDuration(createdAt: string, finishedAt: string): string {
  if (!createdAt && !finishedAt) {
    return ''
  }

  const es = elapsedSeconds(new Date(createdAt), new Date(finishedAt))
  const hours = Math.floor(es / SECONDS_IN_HOUR)
  const minutes = Math.floor((es % SECONDS_IN_HOUR) / SECONDS_IN_MINUTE)
  const seconds = Math.ceil((es % SECONDS_IN_HOUR) % SECONDS_IN_MINUTE)

  return format(hours, minutes, seconds)
}

function elapsedSeconds(createdAt: Date, finishedAt: Date): number {
  return (
    finishedAt.getTime() / MS_IN_SECOND - createdAt.getTime() / MS_IN_SECOND
  )
}

function format(hours: number, minutes: number, seconds: number) {
  return `${hours ? `${hours}h` : ''}${minutes ? `${minutes}m` : ''}${seconds}s`
}
