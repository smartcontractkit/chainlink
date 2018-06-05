export const jobsSelector = state => (
  state
    .jobs
    .currentPage
    .map(id => state.jobs.items[id])
    .filter(j => j)
)

export const jobSpecSelector = (state, jobSpecId) => (
  state.jobs.items[jobSpecId]
)

export const jobRunsSelector = (state, jobSpecId) => {
  const jobSpec = jobSpecSelector(state, jobSpecId)
  const runs = (jobSpec && jobSpec.runs) || []

  return runs
    .map(jobRunId => state.jobRuns.items[jobRunId])
    .filter(r => r)
}

export const latestJobRunsSelector = (state, jobSpecId) => {
  const jobRuns = jobRunsSelector(state, jobSpecId)

  return jobRuns
    .sort((a, b) => {
      const dateA = new Date(a.createdAt)
      const dateB = new Date(b.createdAt)

      return dateA < dateB ? 1 : -1
    })
    .slice(0, 5)
}
