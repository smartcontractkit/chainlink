import 'whatwg-fetch'

const getJobRuns = async (): Promise<IJobRun[]> => {
  const r: Response = await fetch('/api/v1/job_runs')
  return r.json()
}

export { getJobRuns }
