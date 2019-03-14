import 'whatwg-fetch'

const getJobRuns = () => fetch('/api/v1/job_runs').then((r: any) => r.json())

export { getJobRuns }
