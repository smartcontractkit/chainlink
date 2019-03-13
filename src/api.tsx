const getJobRuns = () => fetch('/api/v1/job_runs').then(r => r.json())

export { getJobRuns }
